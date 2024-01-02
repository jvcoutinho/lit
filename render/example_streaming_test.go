package render_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/render"
)

func Stream(r *lit.Request) lit.Response {
	streamContent := strings.NewReader("streaming content")

	return render.Stream(r, streamContent)
}

func Example_streaming() {
	r := lit.NewRouter()
	r.GET("/stream", Stream)

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/stream", nil)
	r.ServeHTTP(res, req)

	fmt.Println(res.Body)
	fmt.Println(res.Header())

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/stream", nil)
	req.Header.Set("If-Match", "tag")
	r.ServeHTTP(res, req)

	fmt.Println(res.Body)

	// Output:
	// streaming content
	// map[Accept-Ranges:[bytes] Content-Length:[17] Content-Type:[text/plain; charset=utf-8]]
	//
}
