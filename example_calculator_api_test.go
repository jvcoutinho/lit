package lit_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/bind"
	"github.com/jvcoutinho/lit/render"
	"github.com/jvcoutinho/lit/validate"
)

type BinaryOperationRequest struct {
	A int `query:"a"`
	B int `query:"b"`
}

// Add returns the sum of two integers.
func Add(r *lit.Request) lit.Response {
	req, err := bind.Query[BinaryOperationRequest](r)
	if err != nil {
		return render.BadRequest(err)
	}

	return render.OK(req.A + req.B)
}

// Subtract returns the subtraction of two integers.
func Subtract(r *lit.Request) lit.Response {
	req, err := bind.Query[BinaryOperationRequest](r)
	if err != nil {
		return render.BadRequest(err)
	}

	return render.OK(req.A - req.B)
}

// Multiply returns the multiplication of two integers.
func Multiply(r *lit.Request) lit.Response {
	req, err := bind.Query[BinaryOperationRequest](r)
	if err != nil {
		return render.BadRequest(err)
	}

	return render.OK(req.A * req.B)
}

// Divide returns the division of two integers, granted the divisor is different from zero.
func Divide(r *lit.Request) lit.Response {
	req, err := bind.Query[BinaryOperationRequest](r)
	if err != nil {
		return render.BadRequest(err)
	}

	if err := validate.Fields(&req,
		validate.NotEqual(&req.B, 0),
	); err != nil {
		return render.BadRequest(err)
	}

	return render.OK(req.A / req.B)
}

func Example_calculatorAPI() {
	r := lit.NewRouter()
	r.GET("/add", Add)
	r.GET("/sub", Subtract)
	r.GET("/mul", Multiply)
	r.GET("/div", Divide)

	request(r, "/add?a=2&b=3")
	request(r, "/sub?a=3&b=3")
	request(r, "/mul?a=3&b=9")
	request(r, "/div?a=6&b=2")

	request(r, "/div?a=6&b=0")
	request(r, "/add?a=2a&b=3")

	// Output:
	// 5
	// 0
	// 27
	// 3
	// {"message":"b should not be equal to 0"}
	// {"message":"a: 2a is not a valid int: invalid syntax"}
}

func request(r *lit.Router, path string) {
	res := httptest.NewRecorder()

	r.ServeHTTP(
		res,
		httptest.NewRequest(http.MethodGet, path, nil),
	)

	fmt.Println(res.Body)
}
