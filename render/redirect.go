package render

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/jvcoutinho/lit"
)

// RedirectResponse sets the Location header to a given target URL and performs a redirect.
type RedirectResponse struct {
	request    *http.Request
	url        string
	statusCode int
}

func (r *RedirectResponse) Write(writer http.ResponseWriter) error {
	parsedURL, err := url.Parse(r.url)
	if err != nil {
		return fmt.Errorf("parsing URL: %w", err)
	}

	http.Redirect(writer, r.request, parsedURL.String(), r.statusCode)

	return nil
}

// Redirect performs an HTTP redirect to a new target URL (set as the Location header).
func Redirect(req *lit.Request, locationURL string, statusCode int) *RedirectResponse {
	return &RedirectResponse{
		req.HTTPRequest(),
		locationURL,
		statusCode,
	}
}

// MovedPermanently performs a permanent redirect with Status Code 301 (Moved Permanently) to a new target URL.
func MovedPermanently(req *lit.Request, targetURL string) *RedirectResponse {
	return Redirect(req, targetURL, http.StatusMovedPermanently)
}

// Found performs a temporary redirect with Status Code 302 (Found) to a new target URL.
func Found(req *lit.Request, targetURL string) *RedirectResponse {
	return Redirect(req, targetURL, http.StatusFound)
}

// SeeOther performs a temporary redirect with Status Code 303 (Found) to a new target URL.
func SeeOther(req *lit.Request, targetURL string) *RedirectResponse {
	return Redirect(req, targetURL, http.StatusSeeOther)
}

// TemporaryRedirect performs a temporary redirect with Status Code 307 (Temporary Redirect) to a new target URL.
func TemporaryRedirect(req *lit.Request, targetURL string) *RedirectResponse {
	return Redirect(req, targetURL, http.StatusTemporaryRedirect)
}

// PermanentRedirect performs a permanent redirect with Status Code 308 (Permanent Redirect) to a new target URL.
func PermanentRedirect(req *lit.Request, targetURL string) *RedirectResponse {
	return Redirect(req, targetURL, http.StatusPermanentRedirect)
}
