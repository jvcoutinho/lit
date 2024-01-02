package lit_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/render"
)

func MethodNotAllowedHandler(_ *lit.Request) lit.Response {
	return render.JSON(http.StatusMethodNotAllowed, "Unsupported method")
}

func Example_customMethodNotAllowedHandler() {
	r := lit.NewRouter()
	r.HandleMethodNotAllowed(MethodNotAllowedHandler)

	r.GET("/", HelloWorld)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	fmt.Println(res.Body, res.Code)

	// Output:
	// {"message":"Unsupported method"} 405
}
