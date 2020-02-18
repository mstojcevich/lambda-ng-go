package upload

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/mstojcevich/lambda-ng-go/database"
	tplt "github.com/mstojcevich/lambda-ng-go/template"
	"github.com/mstojcevich/lambda-ng-go/user"
	"github.com/valyala/fasthttp"
)

var addPasteStmt, _ = database.DB.Prepare(`INSERT INTO pastes (owner, name, content_json, is_code, upload_date) VALUES ($1, $2, $3, $4, $5)`)
var getPasteStmt, _ = database.DB.Prepare(`SELECT content_json FROM pastes WHERE name=$1`)

func init() {
	createPasteTemplate()
	createPasteViewTemplate()
}

func createPasteTemplate() {
	// Create the template
	t, err := template.ParseFiles("html/paste.html", "html/partials/shared_head.html", "html/partials/topbar.html")
	if err != nil {
		panic(err)
	}

	// Render the template into a byte buffer
	var tpl bytes.Buffer
	err = t.Execute(&tpl, tplt.CommonTemplateCtx{NoJS: false})
	if err != nil {
		panic(err)
	}

	// Output the template to a file
	ioutil.WriteFile("html/compiled/paste.html", tpl.Bytes(), 0644)
}

func createPasteViewTemplate() {
	// Create the template
	t, err := template.ParseFiles("html/view_paste.html", "html/partials/shared_head.html", "html/partials/topbar.html")
	if err != nil {
		panic(err)
	}

	// Render the template into a byte buffer
	var tpl bytes.Buffer
	err = t.Execute(&tpl, nil)
	if err != nil {
		panic(err)
	}

	// Output the template to a file
	ioutil.WriteFile("html/compiled/view_paste.html", tpl.Bytes(), 0644)
}

// PastePage renders the paste page HTML
func PastePage(ctx *fasthttp.RequestCtx) {
	ctx.SendFile("html/compiled/paste.html")
}

// GetPasteAPI handles recieving paste content
func GetPasteAPI(ctx *fasthttp.RequestCtx) {
	pasteName := string(ctx.FormValue("name"))
	if len(pasteName) > 0 {
		var pasteContent string

		r := getPasteStmt.QueryRow(pasteName)
		err := r.Scan(&pasteContent)

		if err != nil {
			ctx.Error("Error finding paste", fasthttp.StatusUnprocessableEntity)
			return
		}

		ctx.WriteString(pasteContent)
	} else {
		ctx.Error("No name specified", fasthttp.StatusUnprocessableEntity)
		return
	}
}

// PutPasteAPI handles uploading a paste to Lambda
func PutPasteAPI(ctx *fasthttp.RequestCtx) {
	user, err := user.GetLoggedInUser(ctx)
	if err != nil {
		pasteError(ctx, "You must be logged in to paste", fasthttp.StatusUnauthorized)
		return
	}

	pasteText := string(ctx.FormValue("paste"))
	isCode, _ := strconv.ParseBool(string(ctx.FormValue("is_code")))

	// Generate a name for the paste
	filename, err := genFilename()
	if err != nil {
		pasteError(ctx, "Error generating filename for paste", fasthttp.StatusInternalServerError)
		return
	}

	// Save the paste into the database
	_, err = addPasteStmt.Exec(user.ID, filename, pasteText, isCode, time.Now())
	if err != nil {
		pasteError(ctx, "Error saving paste", fasthttp.StatusInternalServerError)
		return
	}

	// Respond with the paste URL
	result := PasteResponse{URL: filename}
	resultJSON, err := result.MarshalJSON()
	if err != nil {
		log.Println(err)
		ctx.Error("{errors:[\"Failed to create JSON response. Contact an admin\"]}", fasthttp.StatusInternalServerError)
		ctx.SetContentType("text/json")
		return
	}
	ctx.Write(resultJSON)
	ctx.SetContentType("text/json")
}

// pasteError writes out JSON for a failed paste
func pasteError(ctx *fasthttp.RequestCtx, errStr string, statusCode int) {
	result := PasteResponse{Errors: []string{errStr}}
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
