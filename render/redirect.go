package render

import (
	"fmt"
	"net/http"
	"net/url"
)

// RedirectResponse sets the Location header to a given target URL and performs a redirect.
type RedirectResponse struct {
	*HTTPResponse

	url string
}

func (r *RedirectResponse) Write(writer http.ResponseWriter) error {
	parsedURL, err := url.Parse(r.url)
	if err != nil {
		return fmt.Errorf("parsing URL: %w", err)
	}

	r.header.Set("Location", parsedURL.String())

	return r.HTTPResponse.Write(writer)
}

// Redirect performs an HTTP redirect to a new target URL (set as the Location header).
func Redirect(statusCode int, locationURL string) *RedirectResponse {
	return &RedirectResponse{
		NewHTTPResponse(statusCode, nil),
		locationURL,
	}
}

// MovedPermanently performs a permanent redirect with Status Code 301 (Moved Permanently) to a new target URL.
func MovedPermanently(targetURL string) *RedirectResponse {
	return Redirect(http.StatusMovedPermanently, targetURL)
}

// Found performs a temporary redirect with Status Code 302 (Found) to a new target URL.
func Found(targetURL string) *RedirectResponse {
	return Redirect(http.StatusFound, targetURL)
}

// SeeOther performs a temporary redirect with Status Code 303 (Found) to a new target URL.
func SeeOther(targetURL string) *RedirectResponse {
	return Redirect(http.StatusSeeOther, targetURL)
}

// TemporaryRedirect performs a temporary redirect with Status Code 307 (Temporary Redirect) to a new target URL.
func TemporaryRedirect(targetURL string) *RedirectResponse {
	return Redirect(http.StatusTemporaryRedirect, targetURL)
}

// PermanentRedirect performs a permanent redirect with Status Code 308 (Permanent Redirect) to a new target URL.
func PermanentRedirect(targetURL string) *RedirectResponse {
	return Redirect(http.StatusPermanentRedirect, targetURL)
}
