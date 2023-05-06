package lit

// Result of an HTTP request.
//
// Use the result package in order to produce Result implementations.
type Result interface {
	// Write writes the result into the HTTP response managed by ctx.
	Write(ctx *Context)
}
