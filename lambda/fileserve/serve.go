package fileserve

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mstojcevich/lambda-ng-go/config"
	"github.com/mstojcevich/lambda-ng-go/database"
	"github.com/valyala/fasthttp"
	"gopkg.in/kothar/go-backblaze.v0"
)

var B2 *backblaze.B2
var Bucket *backblaze.Bucket

var pasteExistsStmt, _ = database.DB.Prepare(`SELECT exists(SELECT 1 FROM pastes WHERE name=$1)`)

var b2NameStmt, _ = database.DB.Prepare(`
	SELECT CASE
	WHEN in_b2 THEN
		CONCAT(name, '.', extension)
	ELSE
		''
	END
	FROM files WHERE name=$1
`)

func init() {
	create404Template()

	if len(config.BackblazeAccountID) > 0 {
		var err error
		B2, err = backblaze.NewB2(backblaze.Credentials{
			AccountID:      config.BackblazeAccountID,
			ApplicationKey: config.BackblazeAppKey,
		})
		if err != nil {
			fmt.Printf("Backblaze account ID: %s\n", config.BackblazeAccountID)
			fmt.Printf("Backblaze app key: %s\n", config.BackblazeAppKey)
			panic(err)
		}

		Bucket, err = B2.Bucket(config.BackblazeBucket)
		if err != nil {
			panic(err)
		}
	}
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

func show404(ctx *fasthttp.RequestCtx) {
	ctx.SendFile("html/compiled/404.html")
	ctx.SetContentType("text/html")
	ctx.SetStatusCode(404)
}

// viewPastePage renders the paste page HTML
func viewPastePage(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html")
	ctx.SendFile("html/compiled/view_paste.html")
}

// Serve serves an uploaded Lambda file or paste
func Serve(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	path = path[1:len(path)] // exclude leading /

	fasthttp.ServeFileUncompressed(ctx, fmt.Sprintf("files/%s", path))

	foundUpload := ctx.Response.StatusCode() != 404
	if foundUpload {
		return
	}

	var pasteExists bool
	r := pasteExistsStmt.QueryRow(path)
	r.Scan(&pasteExists)

	if pasteExists {
		viewPastePage(ctx)
		return
	}

	// Look for the upload in BackBlaze B2
	dotSplit := strings.Split(path, ".")
	if len(dotSplit) != 2 {
		show404(ctx)
		return
	}
	name := dotSplit[0]

	r = b2NameStmt.QueryRow(name)
	var b2name string
	r.Scan(&b2name)
	if b2name != "" {
		// Grab from B2 then serve
		_, readCloser, err := Bucket.DownloadFileByName(b2name)
		if err != nil {
			fmt.Printf("Error when downloading from B2: %s\n", err)
			show404(ctx)
			return
		}
		defer readCloser.Close()

		outFile, err := os.Create(fmt.Sprintf("files/%s", b2name))
		if err != nil {
			fmt.Printf("Error when downloading from B2: %s\n", err)
			show404(ctx)
			return
		}

		// We could hypothetically copy to FS in a new goroutine while simulataniously
		// streaming to make our response faster, but we don't yet.
		_, err = io.Copy(outFile, readCloser)
		if err != nil {
			fmt.Printf("Error when downloading from B2: %s\n", err)
			show404(ctx)
			return
		}
		outFile.Close()

		ctx.Response.Header.SetContentType("")
		fasthttp.ServeFileUncompressed(ctx, fmt.Sprintf("files/%s", b2name))
		return
	}

	show404(ctx)
	return
}
