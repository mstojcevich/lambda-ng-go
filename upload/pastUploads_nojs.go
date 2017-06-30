package upload

import (
	"bytes"
	"html/template"
	"log"
	"math"
	"strconv"

	"github.com/mstojcevich/lambda-ng-go/user"
	"github.com/valyala/fasthttp"
)

type PastUploadsNoJSContext struct {
	user.AuthedTemplateContext

	PastUploads []PastUpload
	NumPages    int
	PageNum     int
}

var puFuncMap = template.FuncMap{
	"add": func(a int, b int) int {
		return a + b
	},
}

// PastUploadsPage handles viewing the past uploads HTML page
// It is accessable at /upload via GET
func PastUploadsPageNoJS(ctx *fasthttp.RequestCtx) {
	t := template.New("past_uploads.html")

	t.Funcs(puFuncMap)

	// Create the template
	t, err := t.ParseFiles("html/past_uploads.html", "html/partials/shared_head.html", "html/partials/topbar.html")
	if err != nil {
		panic(err)
	}

	user, err := user.GetLoggedInUser(ctx)
	if err != nil { // User isn't logged in, bring them to the login page
		ctx.Redirect("/nojs/login", fasthttp.StatusTemporaryRedirect)
		return
	}

	n := 21
	pageNum, err := strconv.Atoi(string(ctx.FormValue("page")))
	if err != nil {
		pageNum = 1
	}

	var tplContext PastUploadsNoJSContext
	tplContext.NoJS = true
	tplContext.PageNum = pageNum
	tplContext.SignedIn = true
	tplContext.Session = user

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

	if pageNum < 1 {
		pageNum = 1
	}

	rows, err := getUploadsStmt.Query(user.ID, n, (pageNum-1)*n)
	if err != nil {
		log.Println(err)
		pastUploadError(ctx, "Error looking for past uploads", fasthttp.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var pastUpload PastUpload

		err = rows.Scan(&pastUpload.ID, &pastUpload.Name, &pastUpload.LocalName, &pastUpload.Extension, &pastUpload.HasThumbnail)
		if err != nil {
			log.Println(err)
			pastUploadError(ctx, "Error looking for past uploads", fasthttp.StatusInternalServerError)
			return
		}

		tplContext.PastUploads = append(tplContext.PastUploads, pastUpload)
	}

	tplContext.NumPages = numPages

	// Render the template into a byte buffer
	var tpl bytes.Buffer
	err = t.Execute(&tpl, tplContext)
	if err != nil {
		panic(err)
	}

	ctx.Write(tpl.Bytes())
	ctx.SetContentType("text/html")
}
