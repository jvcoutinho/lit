package render_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jvcoutinho/lit/render"
	"github.com/stretchr/testify/require"
)

func TestNoContentResponse(t *testing.T) {
	t.Parallel()

	// Arrange
	var (
		writer   = httptest.NewRecorder()
		response = render.NoContent().
				WithHeader("Key", "Value")
	)

	// Act
	response.Write(writer)

	// Assert
	require.Empty(t, writer.Body)
	require.Equal(t, http.StatusNoContent, writer.Code)
	require.Equal(t, http.Header{"Key": []string{"Value"}}, writer.Header())
}
