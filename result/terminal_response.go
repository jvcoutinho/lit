package result

import (
	"github.com/jvcoutinho/lit"
)

// TerminalResponse writes a single chunk of bytes into the response and closes the request.
type TerminalResponse struct {
	statusCode int
	body       []byte
	header     map[string]string
}

func newTerminalResponse(statusCode int, body []byte) *TerminalResponse {
	return &TerminalResponse{statusCode: statusCode, body: body}
}

func (r *TerminalResponse) AddHeader(key string, value string) {
	if r.header == nil {
		r.header = make(map[string]string)
	}

	r.header[key] = value
}

func (r *TerminalResponse) Write(ctx *lit.Context) {
	responseWriter := ctx.ResponseWriter

	header := responseWriter.Header()
	for key, value := range r.header {
		header.Set(key, value)
	}

	responseWriter.WriteHeader(r.statusCode)

	_, _ = responseWriter.Write(r.body)
}
