package render_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/render"
)

var f *os.File

func ServeFile(r *lit.Request) lit.Response {
	return render.File(r, f.Name())
}

func Example_servingFiles() {
	f = createTempFile()
	defer os.Remove(f.Name())

	r := lit.NewRouter()
	r.GET("/file", ServeFile)

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/file", nil)
	r.ServeHTTP(res, req)

	fmt.Println(res.Body)

	// Output:
	// content
}

func createTempFile() *os.File {
	f, err := os.CreateTemp("", "test_file")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := f.Write([]byte("content")); err != nil {
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	return f
}
