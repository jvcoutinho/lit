package render

import (
	"io"
	"net/http"
	"time"

	"github.com/jvcoutinho/lit"
)

// StreamResponse is a lit.Response that sends separated small chunks of data from a given content.
type StreamResponse struct {
	Request      *lit.Request
	Content      io.ReadSeeker
	FilePath     string
	LastModified time.Time
}

func (r StreamResponse) Write(w http.ResponseWriter) {
	http.ServeContent(w, r.Request.Base(), r.FilePath, r.LastModified, r.Content)
}

// WithFilePath sets the file path property of the stream. If it is set, StreamResponse uses its extension
// to derive the Content-Type header, falling back to the stream content otherwise or if it fails.
func (r StreamResponse) WithFilePath(filePath string) StreamResponse {
	r.FilePath = filePath
	return r
}

// WithLastModified sets the last modified property of the stream. If it is set, StreamResponse includes it in
// the Last-Modified header and, if the request contains an If-Modified-Since header, it uses its value to decide
// whether the content needs to be sent at all.
func (r StreamResponse) WithLastModified(lastModified time.Time) StreamResponse {
	r.LastModified = lastModified
	return r
}

// Stream responds the request with a stream, sending smaller chunks of a possibly large data.
func Stream(r *lit.Request, content io.ReadSeeker) StreamResponse {
	return StreamResponse{r, content, "", time.Time{}}
}
