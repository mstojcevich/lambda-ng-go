package template

// CommonTemplateCtx is the common template context used for almost every page
type CommonTemplateCtx struct {
	// AssetMap is a mapping from asset names to their hashed versions for cache busting
	AssetMap map[string]string
	// NoJS is whether the page should be built without JavaScript
	NoJS bool
}
