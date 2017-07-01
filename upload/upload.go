package upload

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"time"

	bimg "gopkg.in/h2non/bimg.v1"

	"github.com/dchest/uniuri"
	clamd "github.com/dutchcoders/go-clamd"
	"github.com/mstojcevich/lambda-ng-go/config"
	"github.com/mstojcevich/lambda-ng-go/database"
	tplt "github.com/mstojcevich/lambda-ng-go/template"
	"github.com/mstojcevich/lambda-ng-go/user"

	"github.com/valyala/fasthttp"
)

var fileExistsStmt, _ = database.DB.Prepare(`SELECT exists(SELECT 1 FROM files WHERE name=$1) OR exists(SELECT 1 FROM pastes WHERE name=$1)`)
var addFileStmt, _ = database.DB.Prepare(`INSERT INTO files (name, owner, extension, encrypted, local_name, upload_date, has_thumbnail) VALUES ($1, $2, $3, false, $4, $5, false) RETURNING id`)
var addThumbnailStmt, _ = database.DB.Prepare(`INSERT INTO thumbnails (parent_name, width, height, url) VALUES ($1, 128, 128, $2)`)
var markThumbnailStmt, _ = database.DB.Prepare(`UPDATE files SET has_thumbnail=true WHERE id=$1`)

var getFileForDeleteStmt, _ = database.DB.Prepare(`SELECT id,owner,name,extension,has_thumbnail FROM files WHERE name=$1`)
var getPasteForDeleteStmt, _ = database.DB.Prepare(`SELECT id,owner FROM pastes WHERE name=$1`)

var deleteFileStmt, _ = database.DB.Prepare("DELETE FROM files WHERE id=$1")
var deletePasteStmt, _ = database.DB.Prepare("DELETE FROM pastes WHERE id=$1")
var deleteThumbnailsStmt, _ = database.DB.Prepare("DELETE FROM thumbnails WHERE parent_name=$1")

var thumbnailOptions = bimg.Options{
	Width:     128,
	Height:    128,
	Crop:      true,
	Quality:   75,
	Interlace: true,
}

// Set used to quickly check if an extension is allowed
var allowedExtensionsSet = map[string]struct{}{}

var clamav *clamd.Clamd

// uploadTplContext is context used when rendering the upload template
type uploadTplContext struct {
	tplt.CommonTemplateCtx

	AllowedExtensions string
	MaxFilesize       int
}

func init() {
	// Fill the allowedExtensions set
	for _, extension := range config.AllowedFiletypes {
		allowedExtensionsSet[extension] = struct{}{}
		allowedExtensionsSet[strings.ToUpper(extension)] = struct{}{}
	}

	createUploadTemplate()

	if config.ClamAVScanning {
		clamav = clamd.NewClamd(config.ClamSock)
	}
}

func createUploadTemplate() {
	// Create the template
	t, err := template.ParseFiles("html/upload.html", "html/partials/shared_head.html", "html/partials/topbar.html")
	if err != nil {
		panic(err)
	}

	// Render the template into a byte buffer
	var tpl bytes.Buffer
	err = t.Execute(&tpl, uploadTplContext{AllowedExtensions: config.AllowedFiletypesStr, MaxFilesize: config.MaxUploadSize})
	if err != nil {
		panic(err)
	}

	// Output the template to a file
	ioutil.WriteFile("html/compiled/upload.html", tpl.Bytes(), 0644)
}

// Page handles viewing the upload HTML page
// It is accessable at /upload via GET
func Page(ctx *fasthttp.RequestCtx) {
	ctx.SendFile("html/compiled/upload.html")
}

// API handles a file upload API request.
// It requires authentication by either api key or session.
// It is accessable at /api/upload via POST
func API(ctx *fasthttp.RequestCtx) {
	responseURLs := upload(ctx)
	responseURL := ""
	if len(responseURLs) > 0 {
		responseURL = responseURLs[0]
	} else {
		return
	}

	response := Response{URL: responseURL, URLs: responseURLs, Errors: nil}
	responseJSON, err := response.MarshalJSON()

	if err != nil {
		log.Println(err)
		ctx.Error("{errors:[\"Failed to create JSON response. Contact an admin\"]}", fasthttp.StatusInternalServerError)
		ctx.SetContentType("text/json")
		return
	}

	ctx.SetContentType("text/json")
	ctx.Write(responseJSON)
}

func upload(ctx *fasthttp.RequestCtx) (responseURLs []string) {
	user, err := user.GetLoggedInUser(ctx)

	if err != nil {
		uploadError(ctx, "You must be logged in to upload", fasthttp.StatusUnauthorized)
		return
	}

	// Pull the file(s) from the multipart form data with name "file"
	files, err := formFiles(ctx, "file")

	// Check if there was any issue getting the files or no files were uploaded
	if files == nil || err != nil || len(files) == 0 {
		// TODO verbose log error
		uploadError(ctx, "File malformed or unspecified", fasthttp.StatusBadRequest)
		return
	}

	// Go through all of the uploaded files
	for _, file := range files {
		// Grab the extension off of the filename
		dotSplit := strings.Split(file.Filename, ".")
		if len(dotSplit) < 2 {
			uploadError(ctx, "File extension not allowed", fasthttp.StatusUnprocessableEntity)
			return
		}

		extension := dotSplit[len(dotSplit)-1]

		// Check if the extension is allowed
		_, extAllowed := allowedExtensionsSet["."+extension]
		if !extAllowed {
			uploadError(ctx, "File extension not allowed", fasthttp.StatusUnprocessableEntity)
			return
		}

		// Generate a random filename, using the extension from the original upload
		filename, err := genFilename()
		if err != nil {
			log.Println(err)
			uploadError(ctx, "Couldn't find a filename that wasn't in use", fasthttp.StatusInternalServerError)
			return
		}

		fullFilename := filename + "." + extension

		// Register the uplaod in the DB before saving
		fileRow := addFileStmt.QueryRow(filename, user.ID, extension, file.Filename, time.Now())

		// Get the file ID of the new file
		var fileID int
		err = fileRow.Scan(&fileID)
		if err != nil {
			log.Println(err)
			uploadError(ctx, "Error registering upload. Please try again later.", fasthttp.StatusInternalServerError)
			return
		}

		// Full path, including the upload folder
		path := config.UploadDir + fullFilename

		// Read the file into a byte buffer
		f, err := file.Open()
		if err != nil {
			log.Println(err)
			uploadError(ctx, "Error saving upload. Please try again later.", fasthttp.StatusInternalServerError)
			return
		}
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		if err != nil {
			log.Println(err)
			uploadError(ctx, "Error saving upload. Please try again later.", fasthttp.StatusInternalServerError)
			return
		}

		err = ioutil.WriteFile(path, b, 0655)
		if err != nil {
			log.Println(err)
			uploadError(ctx, "Error saving upload. Please try again later.", fasthttp.StatusInternalServerError)
			return
		}

		responseURLs = append(responseURLs, fullFilename)

		// Spawn off a worker to scan with clamav and create a thumbnail
		go func() {
			if config.ClamAVScanning {
				response, err := clamav.ScanStream(bytes.NewBuffer(b), make(chan bool))
				if err != nil {
					fmt.Println("Error scanning for viruses")
					fmt.Println(err)
				} else { // Virus scan successfully ran
					rsp := <-response
					if rsp.Status == clamd.RES_FOUND { // Malware found
						fmt.Println("Malware found in " + filename)
						fmt.Println(rsp)

						// Delete the file
						os.Remove(path)
						deleteFileStmt.Exec(fileID)

						return // End goroutine so thumbnail isn't made
					}
				}
			}

			img := bimg.NewImage(b)
			thumb, err := img.Process(thumbnailOptions)
			if err != nil {
				return
			}

			thumbURL := fmt.Sprintf("thumb_128x128_%s.jpg", filename)
			ioutil.WriteFile(config.UploadDir+thumbURL, thumb, 0655)

			// Register the thumbnail in the database
			_, err = addThumbnailStmt.Exec(filename, "/"+thumbURL)
			if err != nil {
				log.Println(err)
				return
			}

			// Mark that the upload has a thumbnail
			_, err = markThumbnailStmt.Exec(fileID)
			if err != nil {
				log.Println(err)
				return
			}
		}()
	}

	return responseURLs
}

// DeleteAPI handles the API used to delete uploads
func DeleteAPI(ctx *fasthttp.RequestCtx) {
	filename := fmt.Sprintf("%s", ctx.UserValue("name"))

	if len(filename) == 0 {
		deleteError(ctx, "No file specified", fasthttp.StatusUnprocessableEntity)
		return
	}

	user, err := user.GetLoggedInUser(ctx)

	if err != nil {
		deleteError(ctx, "You must be signed in to delete an upload", fasthttp.StatusUnauthorized)
		return
	}

	r := getFileForDeleteStmt.QueryRow(filename)
	var uploadID, uploadOwner int
	var uploadName, uploadExtension string
	var uploadHasThumbnail bool
	err = r.Scan(&uploadID, &uploadOwner, &uploadName, &uploadExtension, &uploadHasThumbnail)

	if err != nil { // Upload probably didn't exist, look for paste
		r = getPasteForDeleteStmt.QueryRow(filename)
		var pasteOwner, pasteID int
		err = r.Scan(&pasteID, &pasteOwner)

		if err != nil {
			deleteError(ctx, "Upload does not exist", fasthttp.StatusUnprocessableEntity)
			return
		}

		// Check to make sure the user has the rights to delete the paste
		if pasteOwner != user.ID {
			deleteError(ctx, "That's not your file to delete", fasthttp.StatusUnauthorized)
			return
		}

		// Delete
		_, err = deletePasteStmt.Exec(pasteID)
		if err != nil {
			log.Println(err)
			deleteError(ctx, "An error occurred while deleting paste", fasthttp.StatusInternalServerError)
			return
		}
	} else { // Handle file delete
		// Check to make sure the user has the rights to delete the file
		if uploadOwner != user.ID {
			deleteError(ctx, "That's not your file to delete", fasthttp.StatusUnauthorized)
			return
		}

		// Delete
		_, err = deleteFileStmt.Exec(uploadID)
		if err != nil {
			log.Println(err)
			deleteError(ctx, "An error occurred while deleting file", fasthttp.StatusInternalServerError)
			return
		}

		// Delete from filesystem
		filePath := config.UploadDir + uploadName + "." + uploadExtension
		err = os.Remove(filePath)
		if err != nil {
			log.Println("Failed to delete " + filePath)
			log.Println(err)
			// Don't error for the user, still successfully removed from DB and probably wasn't in filesystem
		}

		if uploadHasThumbnail {
			// Delete thumbnails from thumbnail DB
			// I don't really care if this fails, so I won't error check
			deleteThumbnailsStmt.Exec(uploadName)

			// Delete thumbnail from FS if exists
			// TODO hardcoded thumbnail URL
			thumbnailPath := config.UploadDir + fmt.Sprintf("thumb_128x128_%s.jpg", filename)
			err = os.Remove(thumbnailPath)
			if err != nil {
				log.Println("Failed to delete " + thumbnailPath)
				log.Println(err)
				// Don't error for the user, still successfully removed from DB and probably wasn't in filesystem
			}
		}
	}

	fmt.Println(user.ID)
}

// genFilename generates a random filename that isn't in use for a file or paste
// It returns an error if it was unable to find a name that didn't exist
func genFilename() (string, error) {
	const maxTries = 10
	const startLength = 3
	const triesPerLengthIncrease = 3

	var name = ""

	tries := 0
	for tries <= maxTries {
		nameLength := startLength + int(tries/triesPerLengthIncrease)

		name = uniuri.NewLen(nameLength)

		var exists bool
		err := fileExistsStmt.QueryRow(name).Scan(&exists)
		if err != nil {
			log.Println(err)
			return "", err
		}

		if !exists {
			return name, nil
		}

		tries = tries + 1
	}

	return name, errors.New("Failed to find a name that wasn't in use")
}

/**
 * Taken from fasthttp, modified to give the entire array of files instead of just the first
 */
func formFiles(ctx *fasthttp.RequestCtx, key string) ([]*multipart.FileHeader, error) {
	mf, err := ctx.MultipartForm()
	if err != nil {
		return nil, err
	}
	if mf.File == nil {
		return nil, err
	}
	fhh := mf.File[key]
	if fhh == nil {
		return nil, fasthttp.ErrMissingFile
	}
	return fhh, nil
}

// uploadError writes out JSON for a failed upload
func uploadError(ctx *fasthttp.RequestCtx, errStr string, statusCode int) {
	result := Response{Errors: []string{errStr}}
	resultJSON, err := result.MarshalJSON()
	if err != nil {
		log.Println(err)
		ctx.Error("{errors:[\"Failed to create JSON response. Contact an admin\"]}", fasthttp.StatusInternalServerError)
		ctx.SetContentType("text/json")
		return
	}
	ctx.Error(string(resultJSON), statusCode)
	ctx.SetContentType("text/json")
}

// deleteError writes out JSON for a failed delete
func deleteError(ctx *fasthttp.RequestCtx, errStr string, statusCode int) {
	result := DeleteResponse{Errors: []string{errStr}}
	resultJSON, err := result.MarshalJSON()
	if err != nil {
		log.Println(err)
		ctx.Error("{errors:[\"Failed to create JSON response. Contact an admin\"]}", fasthttp.StatusInternalServerError)
		ctx.SetContentType("text/json")
		return
	}
	ctx.Error(string(resultJSON), statusCode)
	ctx.SetContentType("text/json")
}
