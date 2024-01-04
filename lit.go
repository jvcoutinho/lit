// Package lit is a fast and expressive HTTP framework.
//
// # The basics
//
// In Lit, an HTTP handler is a function that receives a [*Request] and returns a [Response]. Register new handlers
// using the [*Router.Handle] method.
//
// For instance, the Divide function below is a handler that returns the division of two integers coming from query
// parameters:
//
//	type Request struct {
//		A int `query:"a"`
//		B int `query:"b"`
//	}
//
//	func (r *Request) Validate() []validate.Field {
//		return []validate.Field{
//			validate.NotEqual(&r.B, 0),
//		}
//	}
//
//	func Divide(r *lit.Request) lit.Response {
//		req, err := bind.Query[Request](r)
//		if err != nil {
//			return render.BadRequest(err)
//		}
//
//		return render.OK(req.A / req.B)
//	}
//
// In order to extend a handler's functionality and to be able to reuse the logic in several handlers,
// such as logging or authorization logic, one can use middlewares.
// In Lit, a middleware is a function that receives a [Handler] and returns a [Handler]. Register new middlewares
// using the [*Router.Use] method.
//
// For instance, the AppendRequestID function below is a middleware that assigns an ID to the request and appends it
// to the context:
//
//	type ContextKeyType string
//
//	var RequestIDKey ContextKeyType = "request-id"
//
//	func AppendRequestID(h lit.Handler) lit.Handler {
//		return func(r *lit.Request) lit.Response {
//			var (
//				requestID = uuid.New()
//				ctx       = context.WithValue(r.Context(), RequestIDKey, requestID)
//			)
//
//			r.WithContext(ctx)
//
//			return h(r)
//		}
//	}
//
// It is recommended to use the [Log] and [Recover] middlewares.
//
// Check the [package-level examples] for more use cases.
//
// # Model binding and receiving files
//
// Lit can parse and validate data coming from a request's URI parameters, header, body or query parameters
// to Go structs, including files from multipart form requests.
//
// Check [github.com/jvcoutinho/lit/bind] package.
//
// # Validation
//
// Lit can validate Go structs with generics and compile-time assertions.
//
// Check [github.com/jvcoutinho/lit/validate] package.
//
// # Responding requests, redirecting, serving files and streams
//
// Lit responds requests with implementations of the [Response] interface. Current provided implementations include
// JSON responses, redirections, no content responses, files and streams.
//
// Check [github.com/jvcoutinho/lit/render] package.
//
// # Testing handlers
//
// Handlers can be unit tested in several ways. The simplest and idiomatic form is calling the handler with a crafted
// request and asserting the response:
//
//	type Request struct {
//		A int `query:"a"`
//		B int `query:"b"`
//	}
//
//	func (r *Request) Validate() []validate.Field {
//		return []validate.Field{
//			validate.NotEqual(&r.B, 0),
//		}
//	}
//
//	func Divide(r *lit.Request) lit.Response {
//		req, err := bind.Query[Request](r)
//		if err != nil {
//			return render.BadRequest(err)
//		}
//
//		return render.OK(req.A / req.B)
//	}
//
//	func TestDivide(t *testing.T) {
//		t.Parallel()
//
//		tests := []struct {
//			description string
//			a           int
//			b           int
//			want        lit.Response
//		}{
//			{
//				description: "BEquals0",
//				a:           3,
//				b:           0,
//				want:        render.BadRequest("b should not be equal to 0"),
//			},
//			{
//				description: "Division",
//				a:           6,
//				b:           3,
//				want:        render.OK(2),
//			},
//		}
//
//		for _, test := range tests {
//			test := test
//			t.Run(test.description, func(t *testing.T) {
//				t.Parallel()
//
//				var (
//					path    = fmt.Sprintf("/?a=%d&b=%d", test.a, test.b)
//					request = lit.NewRequest(
//						httptest.NewRequest(http.MethodGet, path, nil),
//					)
//					got  = Divide(request)
//					want = test.want
//				)
//
//				if !reflect.DeepEqual(got, want) {
//					t.Fatalf("got: %v; want: %v", got, want)
//				}
//			})
//		}
//	}
//
// # Testing middlewares
//
// Middlewares can be tested in the same way as handlers (crafting a request and asserting the response of the handler
// after the transformation):
//
//	func ValidateXAPIKeyHeader(h lit.Handler) lit.Handler {
//		return func(r *lit.Request) lit.Response {
//			apiKeyHeader, err := bind.HeaderField[string](r, "X-API-KEY")
//			if err != nil {
//				return render.BadRequest(err)
//			}
//
//			if apiKeyHeader == "" {
//				return render.Unauthorized("API Key must be provided")
//			}
//
//			return h(r)
//		}
//	}
//
//	func TestValidateXAPIKeyHeader(t *testing.T) {
//		t.Parallel()
//
//		testHandler := func(r *lit.Request) lit.Response {
//			return render.NoContent()
//		}
//
//		tests := []struct {
//			description  string
//			apiKeyHeader string
//			want         lit.Response
//		}{
//			{
//				description:  "EmptyHeader",
//				apiKeyHeader: "",
//				want:         render.Unauthorized("API Key must be provided"),
//			},
//			{
//				description:  "ValidAPIKey",
//				apiKeyHeader: "api-key-1",
//				want:         render.NoContent(),
//			},
//		}
//
//		for _, test := range tests {
//			test := test
//			t.Run(test.description, func(t *testing.T) {
//				t.Parallel()
//
//				r := httptest.NewRequest(http.MethodGet, "/", nil)
//				r.Header.Add("X-API-KEY", test.apiKeyHeader)
//
//				var (
//					request = lit.NewRequest(r)
//					got     = ValidateXAPIKeyHeader(testHandler)(request)
//				)
//
//				if !reflect.DeepEqual(got, test.want) {
//					t.Fatalf("got: %v; want: %v", got, test.want)
//				}
//			})
//		}
//	}
//
// [package-level examples]: https://pkg.go.dev/github.com/jvcoutinho/lit#pkg-examples
package lit
