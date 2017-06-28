package index

import (
	"bytes"
	"html/template"

	"github.com/mstojcevich/lambda-ng-go/user"
	"github.com/valyala/fasthttp"
)

func PageNoJS(ctx *fasthttp.RequestCtx) {
	// Create the template
	t, err := template.ParseFiles("html/index.html", "html/partials/shared_head.html", "html/partials/topbar.html")
	if err != nil {
		panic(err)
	}

	// Get the signed in user
	renderCtx := user.AuthedTemplateContext{}
	renderCtx.NoJS = true
	user, err := user.GetLoggedInUser(ctx)
	if err == nil {
		renderCtx.SignedIn = true
		renderCtx.Session = user
	}

	// Render the template into a byte buffer
	var tpl bytes.Buffer
	err = t.Execute(&tpl, renderCtx)
	if err != nil {
		panic(err)
	}

	ctx.Write(tpl.Bytes())
	ctx.SetContentType("text/html")
}
