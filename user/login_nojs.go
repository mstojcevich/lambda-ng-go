package user

import "github.com/mstojcevich/lambda-ng-go/template"

// AuthedTemplateContext is a context to render a template with a user's session
type AuthedTemplateContext struct {
	template.CommonTemplateCtx

	SignedIn bool
	Session  User
}
