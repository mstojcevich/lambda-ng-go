package upload

import (
	"bytes"

	"github.com/mstojcevich/lambda-ng-go/config"
	"github.com/mstojcevich/lambda-ng-go/user"
	"github.com/valyala/fasthttp"
)

type uploadTplContextNoJS struct {
	user.AuthedTemplateContext
	uploadTplContext
	NoJS bool
}

// Page handles viewing the upload HTML page
// It is accessable at /upload via GET
func PageNoJS(ctx *fasthttp.RequestCtx) {
	user, err := user.GetLoggedInUser(ctx)
	if err != nil {
		ctx.Redirect("/nojs/login", fasthttp.StatusTemporaryRedirect)
		return
	}

	// Render the template into a byte buffer
	var tpl bytes.Buffer
	tplCtx := uploadTplContextNoJS{}
	tplCtx.AllowedExtensions = config.AllowedFiletypesStr
	tplCtx.MaxFilesize = config.MaxUploadSize
	tplCtx.NoJS = true
	tplCtx.Session = user
	tplCtx.SignedIn = true
	err = uploadTemplate.Execute(&tpl, tplCtx)
	if err != nil {
		panic(err)
	}

	ctx.Write(tpl.Bytes())
	ctx.SetContentType("text/html")
}

func APINoJS(ctx *fasthttp.RequestCtx) {
	responseURLs := upload(ctx)
	responseURL := ""
	if len(responseURLs) > 0 {
		responseURL = responseURLs[0]
	} else {
		return
	}

	ctx.Redirect("/"+responseURL, fasthttp.StatusFound)
}
