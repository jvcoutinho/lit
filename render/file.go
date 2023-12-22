package render

import (
	"net/http"

	"github.com/jvcoutinho/lit"
)

// FileResponse is a lit.Response that sends small chunks of data of a given file.
//
// If the file does not exist, FileResponse responds the request with [404 Not Found].
//
// [404 Not Found]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/404
type FileResponse struct {
	Request  *lit.Request
	FilePath string
}

func (r FileResponse) Write(w http.ResponseWriter) {
	http.ServeFile(w, r.Request.Request, r.FilePath)
}

// File responds the request with a stream of the contents of a file or directory in path
// (absolute or relative to the current directory).
//
// If the file does not exist, File responds the request with [404 Not Found].
//
// [404 Not Found]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/404
func File(r *lit.Request, path string) FileResponse {
	return FileResponse{r, path}
}
