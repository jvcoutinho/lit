package lit_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/render"
)

// EnableCORS handles a preflight CORS OPTIONS request.
func EnableCORS(_ *lit.Request) lit.Response {
	res := render.NoContent()

	return lit.ResponseFunc(func(w http.ResponseWriter) {
		res.Write(w)

		header := w.Header()
		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Credentials", "false")
		header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
	})
}

func HelloWorld(_ *lit.Request) lit.Response {
	return render.OK("Hello, World!")
}

func Example_corsHandler() {
	r := lit.NewRouter()
	r.HandleOPTIONS(EnableCORS)

	r.GET("/", HelloWorld)

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	fmt.Println(res.Header(), res.Code)

	// Output:
	// map[Access-Control-Allow-Credentials:[false] Access-Control-Allow-Methods:[GET, OPTIONS] Access-Control-Allow-Origin:[*] Allow:[GET, OPTIONS]] 204
}
