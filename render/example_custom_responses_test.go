package render_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/bind"
	"github.com/jvcoutinho/lit/render"
	"github.com/jvcoutinho/lit/validate"
	"gopkg.in/yaml.v3"
)

// YAMLResponse is a lit.Response that prints a YAML formatted-body as response. It sets
// the Content-Type header to "application/x-yaml".
//
// If the response contains a body but its marshalling fails, YAMLResponse responds an Internal Server Error
// with the error message as plain text.
type YAMLResponse struct {
	StatusCode int
	Header     http.Header
	Body       any
}

// YAML responds the request with statusCode and a body marshalled as YAML. Nil body equals empty body.
//
// If body is a string or an error, YAML marshals render.Message with the body assigned to render.Message.Value.
// Otherwise, it marshals the body as is.
func YAML(statusCode int, body any) YAMLResponse {
	switch cast := body.(type) {
	case string:
		return YAMLResponse{statusCode, make(http.Header), render.Message{Message: cast}}
	case error:
		return YAMLResponse{statusCode, make(http.Header), render.Message{Message: cast.Error()}}
	default:
		return YAMLResponse{statusCode, make(http.Header), cast}
	}
}

func (r YAMLResponse) Write(w http.ResponseWriter) {
	responseHeader := w.Header()
	for key := range r.Header {
		responseHeader.Set(key, r.Header.Get(key))
	}

	if r.Body == nil {
		w.WriteHeader(r.StatusCode)
		return
	}

	w.Header().Set("Content-Type", "application/x-yaml")

	if err := r.writeBody(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (r YAMLResponse) writeBody(w http.ResponseWriter) error {
	bodyBytes, err := yaml.Marshal(r.Body)
	if err != nil {
		return err
	}

	w.WriteHeader(r.StatusCode)

	_, err = w.Write(bodyBytes)

	return err
}

type Request struct {
	A int `query:"a"`
	B int `query:"b"`
}

func (r *Request) Validate() []validate.Field {
	return []validate.Field{
		validate.NotEqual(&r.B, 0),
	}
}

func Divide(r *lit.Request) lit.Response {
	req, err := bind.Query[Request](r)
	if err != nil {
		return YAML(http.StatusBadRequest, err)
	}

	return YAML(http.StatusOK, req.A/req.B)
}

func Example_customResponse() {
	r := lit.NewRouter()
	r.GET("/div", Divide)

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/div?a=4&b=2", nil)
	r.ServeHTTP(res, req)
	fmt.Println(res.Body)

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/div?a=2&b=0", nil)
	r.ServeHTTP(res, req)
	fmt.Println(res.Body)

	// Output:
	// 2
	//
	// message: b should not be equal to 0
}
