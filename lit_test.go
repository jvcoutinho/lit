package lit_test

import (
	"fmt"
	"io"
	"net/http/httptest"

	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/bind"
	"github.com/jvcoutinho/lit/render"
	"github.com/jvcoutinho/lit/validate"
)

func Example() {
	r := lit.NewRouter()
	r.Use(lit.Log)

	r.GET("/users", getUsers)

	server := httptest.NewServer(r)
	defer server.Close()

	res, err := server.Client().Get(server.URL + "/users?age=28&name=Ali")
	if err == nil {
		defer res.Body.Close()

		responseBody, _ := io.ReadAll(res.Body)
		fmt.Println(string(responseBody))
	}

	// Output:
	// [{"id":"user_id_1","age":28,"name":"Alicent"},{"id":"user_id_2","age":28,"name":"Alibaba"}]
}

func getUsers(r *lit.Request) lit.Response {
	req, err := bind.Query[getUsersRequest](r)
	if err != nil {
		return render.BadRequest(err)
	}

	// Retrieve users from a database
	users := getFakeUsers(req)

	return render.OK(users)
}

type getUsersRequest struct {
	Age  int    `query:"age"`
	Name string `query:"name"`
}

func (r *getUsersRequest) Validate() []validate.Field {
	return []validate.Field{
		validate.Greater(&r.Age, 0),
		validate.NotEmpty(&r.Name),
	}
}

type user struct {
	ID   string `json:"id"`
	Age  int    `json:"age"`
	Name string `json:"name"`
}

func getFakeUsers(r getUsersRequest) []user {
	return []user{
		{
			ID:   "user_id_1",
			Age:  r.Age,
			Name: "Alicent",
		},
		{
			ID:   "user_id_2",
			Age:  r.Age,
			Name: "Alibaba",
		},
	}
}
