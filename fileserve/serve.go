package fileserve

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"

	"github.com/mstojcevich/lambda-ng-go/database"
	"github.com/mstojcevich/lambda-ng-go/upload"
	"github.com/valyala/fasthttp"
)

var pasteExistsStmt, _ = database.DB.Prepare(`SELECT exists(SELECT 1 FROM pastes WHERE name=$1)`)

func init() {
	create404Template()
}

func create404Template() {
	// Create the template
	t, err := template.ParseFiles("html/404.html", "html/partials/shared_head.html", "html/partials/topbar.html")
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
	ioutil.WriteFile("html/compiled/404.html", tpl.Bytes(), 0644)
}

func Show404(ctx *fasthttp.RequestCtx) {
	ctx.SendFile("html/compiled/404.html")
}

// Serve serves an uploaded Lambda file or paste
func Serve(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())

	fasthttp.ServeFileUncompressed(ctx, fmt.Sprintf("files/%s", path))

	foundUpload := ctx.Response.StatusCode() != 404

	if !foundUpload { // look for paste
		r := pasteExistsStmt.QueryRow(path[1:len(path)])

		var pasteExists bool
		r.Scan(&pasteExists)

		if pasteExists {
			upload.ViewPastePage(ctx)
		} else {
			Show404(ctx)
		}
	}
}
