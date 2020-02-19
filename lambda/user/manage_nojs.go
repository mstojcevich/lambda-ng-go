package user

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/valyala/fasthttp"
)

func init() {
	createUserManageTemplate()
}

func ManagePageNoJS(ctx *fasthttp.RequestCtx) {
	// Create the template
	t, err := template.ParseFiles("html/user_manage.html", "html/partials/shared_head.html", "html/partials/topbar.html")
	if err != nil {
		panic(err)
	}

	// Get the signed in user
	renderCtx := AuthedTemplateContext{}
	renderCtx.NoJS = true
	user, err := GetLoggedInUser(ctx)
	if err != nil {
		fmt.Println(err)
	}
	if user != nil && err == nil {
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
