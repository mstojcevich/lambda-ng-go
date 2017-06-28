package index

import (
	"bytes"
	"html/template"
	"io/ioutil"

	"github.com/valyala/fasthttp"

	tplt "github.com/mstojcevich/lambda-ng-go/template"
)

func init() {
	createIndexTemplate()
}

func createIndexTemplate() {
	// Create the template
	t, err := template.ParseFiles("html/index.html", "html/partials/shared_head.html", "html/partials/topbar.html")
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
	ioutil.WriteFile("html/compiled/index.html", tpl.Bytes(), 0644)
}

func Page(ctx *fasthttp.RequestCtx) {
	ctx.SendFile("html/compiled/index.html")
}
