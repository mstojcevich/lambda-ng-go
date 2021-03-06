package index

import (
	"bytes"
	"fmt"

	"github.com/mstojcevich/lambda-ng-go/assetmap"
	"github.com/mstojcevich/lambda-ng-go/user"
	"github.com/valyala/fasthttp"
)

func PageNoJS(ctx *fasthttp.RequestCtx) {
	// Get the signed in user
	renderCtx := user.AuthedTemplateContext{}
	renderCtx.AssetMap = assetmap.Assets.Map
	renderCtx.NoJS = true
	user, err := user.GetLoggedInUser(ctx)
	if err != nil {
		fmt.Println(err)
	}
	if err == nil && user != nil {
		renderCtx.SignedIn = true
		renderCtx.Session = user
	}

	// Render the template into a byte buffer
	var tpl bytes.Buffer
	err = indexTemplate.Execute(&tpl, renderCtx)
	if err != nil {
		panic(err)
	}

	ctx.Write(tpl.Bytes())
	ctx.SetContentType("text/html")
}
