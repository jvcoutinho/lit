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

// Even validates if target is an even number.
func Even(target *int) validate.Field {
	return validate.Field{
		Valid:   target != nil && *target%2 == 0,
		Message: "{0} should be an odd number",
		Fields:  []any{target},
	}
}

// MultipleField validates if target is multiple of field.
func MultipleField(target *int, field *int) validate.Field {
	return validate.Field{
		Valid:   target != nil && field != nil && *target%*field == 0,
		Message: "{1} should be divisible by {0}",
		Fields:  []any{target, field},
	}
}

type EvenNumbersBetweenRequest struct {
	A int `query:"a"`
	B int `query:"b"`
}

func (r *EvenNumbersBetweenRequest) Validate() []validate.Field {
	return []validate.Field{
		Even(&r.A),
		MultipleField(&r.B, &r.A),
	}
}

// EvenNumbersBetween compute the even numbers between two numbers, given the first is even and
// the second is multiple of the first.
func EvenNumbersBetween(r *lit.Request) lit.Response {
	req, err := bind.Query[EvenNumbersBetweenRequest](r)
	if err != nil {
		return render.BadRequest(err)
	}

	evenNumbersBetween := make([]int, 0)
	for i := req.A; i <= req.B; i += 2 {
		evenNumbersBetween = append(evenNumbersBetween, i)
	}

	return render.OK(evenNumbersBetween)
}

func Example_customValidations() {
	r := lit.NewRouter()
	r.GET("/", EvenNumbersBetween)

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?a=4&b=20", nil)
	r.ServeHTTP(res, req)
	fmt.Println(res.Body)

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/?a=2&b=3", nil)
	r.ServeHTTP(res, req)
	fmt.Println(res.Body)

	// Output:
	// [4,6,8,10,12,14,16,18,20]
	// {"message":"a should be divisible by b"}
}
