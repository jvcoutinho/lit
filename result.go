package lit

// Result of an HTTP request.
//
// See the lit/render package.
type Result interface {
	// Render writes this into the HTTP response managed by ctx.
	Render(ctx *Context)
}
