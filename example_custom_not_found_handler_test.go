package lit_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/render"
)

func NotFoundHandler(_ *lit.Request) lit.Response {
	return render.NotFound("Not found. Try again later.")
}

func Example_customNotFoundHandler() {
	r := lit.NewRouter()
	r.HandleNotFound(NotFoundHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	fmt.Println(res.Body, res.Code)

	// Output:
	// {"message":"Not found. Try again later."} 404
}
