package user

import (
	"bytes"
	"html/template"
	"io/ioutil"

	"github.com/mstojcevich/lambda-ng-go/assetmap"
	tplt "github.com/mstojcevich/lambda-ng-go/template"
	"github.com/valyala/fasthttp"
)

func init() {
	createUserManageTemplate()
}

func createUserManageTemplate() {
	// Create the template
	t, err := template.ParseFiles("html/user_manage.html", "html/partials/shared_head.html", "html/partials/topbar.html")
	if err != nil {
		panic(err)
	}

	// Render the template into a byte buffer
	var tpl bytes.Buffer
	err = t.Execute(&tpl, tplt.CommonTemplateCtx{AssetMap: assetmap.Assets.Map, NoJS: false})
	if err != nil {
		panic(err)
	}

	// Output the template to a file
	ioutil.WriteFile("html/compiled/user_manage.html", tpl.Bytes(), 0644)
}

func ManagePage(ctx *fasthttp.RequestCtx) {
	ctx.SendFile("html/compiled/user_manage.html")
}
