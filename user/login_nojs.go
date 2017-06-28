package user

import "github.com/mstojcevich/lambda-ng-go/template"
import "github.com/valyala/fasthttp"
import "github.com/mstojcevich/lambda-ng-go/user/session"

// AuthedTemplateContext is a context to render a template with a user's session
type AuthedTemplateContext struct {
	template.CommonTemplateCtx

	SignedIn bool
	Session  User
}

func LogoutNoJS(ctx *fasthttp.RequestCtx) {
	session.Sessions.DestroyFasthttp(ctx)
	ctx.Redirect("/nojs/", fasthttp.StatusFound)
}
