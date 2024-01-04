package render_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/render"
	"github.com/stretchr/testify/require"
)

func TestStreamResponse(t *testing.T) {
	t.Parallel()

	// Arrange
	var (
		writer   = httptest.NewRecorder()
		reader   = bytes.NewReader([]byte("content"))
		request  = lit.NewRequest(httptest.NewRequest(http.MethodGet, "/stream", nil))
		response = render.Stream(request, reader).
				WithFilePath("./stream_test.go").
				WithLastModified(time.Date(2023, 20, 10, 10, 10, 10, 20, time.UTC))
	)

	// Act
	response.Write(writer)

	// Assert
	require.Equal(t, http.StatusOK, writer.Code)
	require.Equal(t, "content", writer.Body.String())
}
