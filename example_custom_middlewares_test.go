package lit_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/jvcoutinho/lit/validate"

	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/bind"
	"github.com/jvcoutinho/lit/render"
)

type GetUserNameRequest struct {
	UserID string `uri:"user_id"`
}

func (r *GetUserNameRequest) Validate() []validate.Field {
	return []validate.Field{
		validate.UUID(&r.UserID),
	}
}

type GetUserNameResponse struct {
	Name string `json:"name"`
}

// GetUserName gets the name of an identified user.
func GetUserName(r *lit.Request) lit.Response {
	_, err := bind.URIParameters[GetUserNameRequest](r)
	if err != nil {
		return render.BadRequest(err)
	}

	// getting user name...
	res := GetUserNameResponse{"John"}

	return render.OK(res)
}

type PatchUserNameRequest struct {
	UserID string `uri:"user_id"`
	Name   string `json:"name"`
}

func (r *PatchUserNameRequest) Validate() []validate.Field {
	return []validate.Field{
		validate.UUID(&r.UserID),
		validate.MinLength(&r.Name, 3),
	}
}

// PatchUserName patches the name of an identified user.
func PatchUserName(r *lit.Request) lit.Response {
	_, err := bind.Request[PatchUserNameRequest](r)
	if err != nil {
		return render.BadRequest(err)
	}

	// patching user name...

	return render.NoContent()
}

func Example_customMiddlewares() {
	r := lit.NewRouter()
	r.Use(ComputeLatency)

	r.GET("/users/:user_id/name", GetUserName)
	r.PATCH("/users/:user_id/name", PatchUserName,
		ValidateXAPIKeyHeader)

	res := httptest.NewRecorder()
	r.ServeHTTP(
		res,
		httptest.NewRequest(http.MethodGet, "/users/19fb2f66-f335-47ef-a1ca-1d02d1a117c8/name", nil),
	)
	fmt.Println(res.Body, res.Code)

	res = httptest.NewRecorder()
	r.ServeHTTP(
		res,
		httptest.NewRequest(http.MethodPatch, "/users/19fb2f66-f335-47ef-a1ca-1d02d1a117c8/name",
			strings.NewReader(`{"name":"John"}`),
		),
	)
	fmt.Println(res.Body, res.Code)

	// Output:
	// Latency of the request is 0s
	// {"name":"John"} 200
	// Latency of the request is 0s
	// {"message":"X-API-Key must be provided"} 401
}

// ComputeLatency computes the latency of a request.
func ComputeLatency(h lit.Handler) lit.Handler {
	return func(r *lit.Request) lit.Response {
		var (
			startTime = time.Now()
			res       = h(r)
			endTime   = time.Now()
		)

		fmt.Printf("Latency of the request is %s\n", endTime.Sub(startTime).Round(time.Second))

		return res
	}
}

// ValidateXAPIKeyHeader validates if the X-API-Key header is not empty.
func ValidateXAPIKeyHeader(h lit.Handler) lit.Handler {
	return func(r *lit.Request) lit.Response {
		apiKeyHeader, err := bind.HeaderField[string](r, "X-API-KEY")
		if err != nil {
			return render.BadRequest(err)
		}

		if apiKeyHeader == "" {
			return render.Unauthorized("X-API-Key must be provided")
		}

		fmt.Printf("Authorized request for %s\n", apiKeyHeader)

		return h(r)
	}
}
