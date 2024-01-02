package render_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/render"
)

func DeprecatedHelloWorld(r *lit.Request) lit.Response {
	return render.PermanentRedirect(r, "/new")
}

func HelloWorld(_ *lit.Request) lit.Response {
	return render.OK("Hello, World!")
}

func Example_redirecting() {
	r := lit.NewRouter()
	r.GET("/", DeprecatedHelloWorld)
	r.GET("/new", HelloWorld)

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(res, req)

	fmt.Println(res.Header(), res.Code)
	fmt.Println(res.Body)

	// Output:
	// map[Content-Type:[text/html; charset=utf-8] Location:[/new]] 308
	// <a href="/new">Permanent Redirect</a>.
}
