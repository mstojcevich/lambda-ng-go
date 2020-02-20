package index

import (
	"bytes"
	"html/template"
	"io/ioutil"

	"github.com/valyala/fasthttp"

	"github.com/mstojcevich/lambda-ng-go/assetmap"
	tplt "github.com/mstojcevich/lambda-ng-go/template"
)

var indexTemplate *template.Template

func init() {
	createIndexTemplate()
}

func createIndexTemplate() {
	// Create the template
	var err error
	indexTemplate, err = template.ParseFiles("html/index.html", "html/partials/shared_head.html", "html/partials/topbar.html")
	if err != nil {
		panic(err)
	}

	// Render the template into a byte buffer
	var tpl bytes.Buffer
	err = indexTemplate.Execute(&tpl, tplt.CommonTemplateCtx{AssetMap: assetmap.Assets.Map, NoJS: false})
	if err != nil {
		panic(err)
	}

	// Output the template to a file
	ioutil.WriteFile("html/compiled/index.html", tpl.Bytes(), 0644)
}

func Page(ctx *fasthttp.RequestCtx) {
	ctx.SendFile("html/compiled/index.html")
}
