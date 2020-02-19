package upload

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"strconv"

	"github.com/mstojcevich/lambda-ng-go/database"
	tplt "github.com/mstojcevich/lambda-ng-go/template"
	"github.com/mstojcevich/lambda-ng-go/user"
	"github.com/valyala/fasthttp"
)

var numUploadsStmt, _ = database.DB.Prepare(`SELECT count(id) FROM files WHERE owner=$1`)
var getUploadsStmt, _ = database.DB.Prepare(`SELECT id,name,local_name,extension,has_thumbnail FROM files WHERE owner=$1 ORDER BY id DESC LIMIT $2 OFFSET $3`)

var pastUploadTemplate *template.Template

func init() {
	createPastUploadsTemplate()
}

func createPastUploadsTemplate() {
	// Create the template
	pastUploadTemplate = template.New("past_uploads.html")

	pastUploadTemplate.Funcs(puFuncMap)

	var err error
	_, err = pastUploadTemplate.ParseFiles("html/past_uploads.html", "html/partials/shared_head.html", "html/partials/topbar.html")
	if err != nil {
		panic(err)
	}

	// Render the template into a byte buffer
	var tpl bytes.Buffer
	err = pastUploadTemplate.Execute(&tpl, tplt.CommonTemplateCtx{NoJS: false})
	if err != nil {
		panic(err)
	}

	// Output the template to a file
	ioutil.WriteFile("html/compiled/past_uploads.html", tpl.Bytes(), 0644)
}

// PastUploadsPage handles viewing the past uploads HTML page
// It is accessable at /upload via GET
func PastUploadsPage(ctx *fasthttp.RequestCtx) {
	ctx.SendFile("html/compiled/past_uploads.html")
}

// PastUploadsAPI enumerates a user's past uploads
func PastUploadsAPI(ctx *fasthttp.RequestCtx) {
	user, err := user.GetLoggedInUser(ctx)
	if err != nil {
		fmt.Println(err)
	}
	if user == nil || err != nil {
		pastUploadError(ctx, "You must be signed in to see past uploads", fasthttp.StatusUnauthorized)
		return
	}

	n, err := strconv.Atoi(string(ctx.FormValue("n")))
	if err != nil {
		pastUploadError(ctx, "Error parsing request. \"n\" should be an integer.", fasthttp.StatusBadRequest)
		return
	}

	pageNum, err := strconv.Atoi(string(ctx.FormValue("page")))
	if err != nil {
		pastUploadError(ctx, "Error parsing request. \"page\" should be an integer.", fasthttp.StatusBadRequest)
		return
	}

	// Limit n between 1 and 50
	if n > 50 {
		n = 50
	}
	if n < 1 {
		n = 1
	}

	// Limit pageNum between 1 and 50
	if pageNum > 15 {
		pageNum = 15
	}
	if pageNum < 1 {
		pageNum = 1
	}

	// Get the number of uploads by the user to calculate the total number of pages
	var numUploads int
	numUploadRow := numUploadsStmt.QueryRow(user.ID)
	err = numUploadRow.Scan(&numUploads)
	if err != nil {
		log.Println(err)
		pastUploadError(ctx, "Error looking for past uploads", fasthttp.StatusInternalServerError)
		return
	}
	numPages := int(math.Min(15, (math.Ceil(float64(numUploads) / float64(n)))))

	rows, err := getUploadsStmt.Query(user.ID, n, (pageNum-1)*n)
	if err != nil {
		log.Println(err)
		pastUploadError(ctx, "Error looking for past uploads", fasthttp.StatusInternalServerError)
		return
	}

	response := PastUploads{NumPages: numPages}
	for rows.Next() {
		var pastUpload PastUpload

		err = rows.Scan(&pastUpload.ID, &pastUpload.Name, &pastUpload.LocalName, &pastUpload.Extension, &pastUpload.HasThumbnail)
		if err != nil {
			log.Println(err)
			pastUploadError(ctx, "Error looking for past uploads", fasthttp.StatusInternalServerError)
			return
		}

		response.Files = append(response.Files, pastUpload)
	}

	resultJSON, err := response.MarshalJSON()
	if err != nil {
		log.Println(err)
		ctx.Error("{errors:[\"Failed to create JSON response. Contact an admin\"]}", fasthttp.StatusInternalServerError)
		ctx.SetContentType("text/json")
		return
	}
	ctx.Write(resultJSON)
	ctx.SetContentType("text/json")
}

// pastUplaodError writes out JSON for a failed listing of past uploads
func pastUploadError(ctx *fasthttp.RequestCtx, errStr string, statusCode int) {
	result := PastUploads{Errors: []string{errStr}}
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

// GenericImageByExtension serves a generic image based on a file extension
func GenericImageByExtension(ctx *fasthttp.RequestCtx) {
	extension := fmt.Sprintf("%s", ctx.UserValue("extension"))

	var genericImg string

	switch extension {
	case "png", "jpg", "jpeg", "svg", "tiff", "webp":
		genericImg = "static/img/generic/image.svg"
	case "mp4", "webm", "avi", "m4v":
		genericImg = "static/img/generic/video.svg"
	case "opus", "ogg", "m4a", "mp3":
		genericImg = "static/img/generic/audio.svg"
	default:
		genericImg = "static/img/generic/generic.svg"
	}

	ctx.SendFile(genericImg)
}
