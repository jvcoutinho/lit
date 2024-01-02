package validate_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/bind"
	"github.com/jvcoutinho/lit/render"
	"github.com/jvcoutinho/lit/validate"
)

type DivideRequest struct {
	A int `query:"a"`
	B int `query:"b"`
}

func (r *DivideRequest) Validate() []validate.Field {
	notZeroValidation := validate.NotEqual(&r.B, 0)
	notZeroValidation.Message = "{0} should not be zero: invalid division"

	return []validate.Field{
		notZeroValidation,
	}
}

// Divide returns the division of two integers, given the second is not zero.
func Divide(r *lit.Request) lit.Response {
	req, err := bind.Query[DivideRequest](r)
	if err != nil {
		return render.BadRequest(err)
	}

	return render.OK(req.A / req.B)
}

func Example_customMessage() {
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
	// {"message":"b should not be zero: invalid division"}
}
