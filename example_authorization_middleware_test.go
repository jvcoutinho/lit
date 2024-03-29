package lit_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/bind"
	"github.com/jvcoutinho/lit/render"
	"github.com/jvcoutinho/lit/validate"
)

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

func Example_authorizationMiddlewares() {
	r := lit.NewRouter()

	r.GET("/users/:user_id/name", GetUserName)
	r.PATCH("/users/:user_id/name", PatchUserName, ValidateXAPIKeyHeader)

	requestAPIKey(r, http.MethodGet, "/users/19fb2f66-f335-47ef-a1ca-1d02d1a117c8/name", "", nil)
	requestAPIKey(r, http.MethodPatch, "/users/19fb2f66-f335-47ef-a1ca-1d02d1a117c8/name",
		"", strings.NewReader(`{"name":"John"}`))
	requestAPIKey(r, http.MethodPatch, "/users/19fb2f66-f335-47ef-a1ca-1d02d1a117c8/name",
		"api-key-1", strings.NewReader(`{"name":"John"}`))

	// Output:
	// {"name":"John"} 200
	// {"message":"X-API-Key must be provided"} 401
	// Authorized request for api-key-1
	//  204
}

func requestAPIKey(r *lit.Router, method, path, header string, body io.Reader) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, body)
	req.Header.Add("X-API-KEY", header)

	r.ServeHTTP(res, req)

	fmt.Println(res.Body, res.Code)
}
