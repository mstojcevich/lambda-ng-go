package user

import (
	"bytes"
	"html/template"

	tmpl "github.com/mstojcevich/lambda-ng-go/template"
	"github.com/mstojcevich/lambda-ng-go/user/session"
	"github.com/valyala/fasthttp"
)

// AuthedTemplateContext is a context to render a template with a user's session
type AuthedTemplateContext struct {
	tmpl.CommonTemplateCtx

	SignedIn bool
	Session  User
}

func LogoutNoJS(ctx *fasthttp.RequestCtx) {
	session.Sessions.DestroyFasthttp(ctx)
	ctx.Redirect("/nojs/", fasthttp.StatusFound)
}

func LoginAPINoJS(ctx *fasthttp.RequestCtx) {
	LoginAPI(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		return
	}

	ctx.Redirect("/", fasthttp.StatusFound)
}

func LoginPageNoJS(ctx *fasthttp.RequestCtx) {
	// Create the template
	t, err := template.ParseFiles("html/login.html", "html/partials/shared_head.html")
	if err != nil {
		panic(err)
	}

	// Get the signed in user
	renderCtx := AuthedTemplateContext{}
	renderCtx.NoJS = true
	user, err := GetLoggedInUser(ctx)
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
