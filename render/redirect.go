package render

import (
	"net/http"

	"github.com/jvcoutinho/lit"
)

// RedirectResponse is a lit.Response that performs redirects.
type RedirectResponse struct {
	StatusCode  int
	Request     *lit.Request
	LocationURL string
}

func (r RedirectResponse) Write(w http.ResponseWriter) {
	http.Redirect(w, r.Request.Request, r.LocationURL, r.StatusCode)
}

// Redirect responds the request with a redirection status code and a target URL (absolute or relative to the
// request path) in the Location header.
//
//   - If permanent is true and preserveMethod is false, it responds with [301 Moved Permanently].
//   - If permanent is false and preserveMethod is false, it responds with [302 Found].
//   - If permanent is false and preserveMethod is true, it responds with [307 Temporary Redirect].
//   - If permanent is true and preserveMethod is true, it responds with [308 Permanent Redirect].
//
// [301 Moved Permanently]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/301
// [302 Found]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/302
// [307 Temporary Redirect]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/307
// [308 Permanent Redirect]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/308
func Redirect(r *lit.Request, locationURL string, permanent, preserveMethod bool) RedirectResponse {
	if !permanent {
		if !preserveMethod {
			return Found(r, locationURL)
		}

		return TemporaryRedirect(r, locationURL)
	}

	if !preserveMethod {
		return MovedPermanently(r, locationURL)
	}

	return PermanentRedirect(r, locationURL)
}

// MovedPermanently responds the request with [301 Moved Permanently] and a target URL (absolute or relative to the
// request path) in the Location header.
//
// [301 Moved Permanently]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/301
func MovedPermanently(r *lit.Request, locationURL string) RedirectResponse {
	return RedirectResponse{http.StatusMovedPermanently, r, locationURL}
}

// Found responds the request with [302 Found] and a target URL (absolute or relative to the
// request path) in the Location header.
//
// [302 Found]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/302
func Found(r *lit.Request, locationURL string) RedirectResponse {
	return RedirectResponse{http.StatusFound, r, locationURL}
}

// TemporaryRedirect responds the request with [307 Temporary Redirect] and a target URL (absolute or relative to the
// request path) in the Location header.
//
// [307 Temporary Redirect]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/307
func TemporaryRedirect(r *lit.Request, locationURL string) RedirectResponse {
	return RedirectResponse{http.StatusTemporaryRedirect, r, locationURL}
}

// PermanentRedirect responds the request with [307 Permanent Redirect] and a target URL (absolute or relative to the
// request path) in the Location header.
//
// [308 Permanent Redirect]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/308
func PermanentRedirect(r *lit.Request, locationURL string) RedirectResponse {
	return RedirectResponse{http.StatusPermanentRedirect, r, locationURL}
}

// SeeOther responds the request with [303 See Other] and a target URL (absolute or relative to the
// request path) in the Location header.
//
// [303 See Other]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/303
func SeeOther(r *lit.Request, locationURL string) RedirectResponse {
	return RedirectResponse{http.StatusSeeOther, r, locationURL}
}
