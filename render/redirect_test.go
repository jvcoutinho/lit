package render_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/render"
	"github.com/stretchr/testify/require"
)

func TestRedirectResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description        string
		response           func(r *lit.Request, locationURL string) render.RedirectResponse
		expectedStatusCode int
	}{
		{
			description: "MovedPermanently_ShouldRespond301",
			response: func(r *lit.Request, locationURL string) render.RedirectResponse {
				return render.MovedPermanently(r, locationURL)
			},
			expectedStatusCode: http.StatusMovedPermanently,
		},
		{
			description: "Redirect_PermanentAndNotPreserveMethod_ShouldRespond301",
			response: func(r *lit.Request, locationURL string) render.RedirectResponse {
				return render.Redirect(r, locationURL, true, false)
			},
			expectedStatusCode: http.StatusMovedPermanently,
		},
		{
			description: "PermanentRedirect_ShouldRespond308",
			response: func(r *lit.Request, locationURL string) render.RedirectResponse {
				return render.PermanentRedirect(r, locationURL)
			},
			expectedStatusCode: http.StatusPermanentRedirect,
		},
		{
			description: "Redirect_PermanentAndPreserveMethod_ShouldRespond308",
			response: func(r *lit.Request, locationURL string) render.RedirectResponse {
				return render.Redirect(r, locationURL, true, true)
			},
			expectedStatusCode: http.StatusPermanentRedirect,
		},
		{
			description: "Found_ShouldRespond302",
			response: func(r *lit.Request, locationURL string) render.RedirectResponse {
				return render.Found(r, locationURL)
			},
			expectedStatusCode: http.StatusFound,
		},
		{
			description: "Redirect_NotPermanentAndNotPreserveMethod_ShouldRespond302",
			response: func(r *lit.Request, locationURL string) render.RedirectResponse {
				return render.Redirect(r, locationURL, false, false)
			},
			expectedStatusCode: http.StatusFound,
		},
		{
			description: "TemporaryRedirect_ShouldRespond307",
			response: func(r *lit.Request, locationURL string) render.RedirectResponse {
				return render.TemporaryRedirect(r, locationURL)
			},
			expectedStatusCode: http.StatusTemporaryRedirect,
		},
		{
			description: "Redirect_NotPermanentAndPreserveMethod_ShouldRespond307",
			response: func(r *lit.Request, locationURL string) render.RedirectResponse {
				return render.Redirect(r, locationURL, false, true)
			},
			expectedStatusCode: http.StatusTemporaryRedirect,
		},
		{
			description: "SeeOther_ShouldRespond303",
			response: func(r *lit.Request, locationURL string) render.RedirectResponse {
				return render.SeeOther(r, locationURL)
			},
			expectedStatusCode: http.StatusSeeOther,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			var (
				writer  = httptest.NewRecorder()
				request = lit.NewRequest(
					httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString("body")),
					nil,
				)
				responses = []render.RedirectResponse{
					test.response(request, "https://redirect-target.com"),
					test.response(request, "/redirect-target"),
				}
			)

			for _, response := range responses {
				// Act
				response.Write(writer)

				// Assert
				require.Equal(t, test.expectedStatusCode, writer.Code)
				require.Equal(t, response.LocationURL, writer.Header().Get("Location"))
			}
		})
	}
}
