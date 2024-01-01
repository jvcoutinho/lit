// Package lit is a fast and expressive HTTP framework.
//
// # The basics
//
// In Lit, an HTTP handler is a function that receives a *lit.Request and returns a lit.Response. Register new handlers
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
// In Lit, a middleware is a function that receives a Handler and returns a Handler. Register new middlewares
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
// It is highly recommended to always use the [lit.Log] and [lit.Recover] middlewares, in this order.
//
// Check the other examples to explore the capabilities of Lit.
//
// # Model binding and validation
//
// Lit can parse and validate string data coming from a request's URI parameters, header, body or query parameters
// to Go structs with the [lit/bind] and [lit/validate] packages' functions. They rely on parametric polymorphism
// (generics) for their functioning.
//
// There are two ways a Go struct can be validated:
//
//   - If the target of a binding function is a struct that implements [validate.Validatable] with a pointer receiver,
//     the function automatically validates it.
//   - Otherwise, the struct can be explicitly validated by calling the [validate.Fields] function.
//
// For instance, below there is a handler that patches the name, e-mail and birth year of a user:
//
//	const minBirthYear = 1998
//
//	type UpdateUserRequest struct {
//		ID        string `uri:"user_id"`
//		Name      string `json:"name"`
//		Email     string `json:"email"`
//		BirthYear int    `json:"birth_year"`
//	}
//
//	func UpdateUser(r *lit.Request) lit.Response {
//		req, err := bind.Request[UpdateUserRequest](r)
//		if err != nil {
//			return render.BadRequest(err)
//		}
//
//		if err := validate.Fields(&req,
//			validate.UUID(&req.ID),
//			validate.NotEmpty(&req.Name),
//			validate.Email(&req.Email),
//			validate.GreaterOrEqual(&req.BirthYear, minBirthYear),
//		); err != nil {
//			return render.BadRequest(err)
//		}
//
//		// Patching user in database...
//
//		return render.NoContent()
//	}
//
// The [lit/validate] package contains several common validations, but creating a custom validation compatible with Lit
// is as easy as creating a new [validate.Field] object!
//
// For example, below there is a new validation that checks if a time field is after another time field. If the
// validation fails, Lit produces a [*validate.Error] with the message "since should be after until".
//
//	type GetBooksFromReleaseDate struct {
//		Since time.Time `json:"since"`
//		Until time.Time `json:"until"`
//	}
//
//	func (r *GetBooksFromReleaseDate) Validate() []validate.Field {
//		return []validate.Field{
//			AfterField(&r.Until, &r.Since),
//		}
//	}
//
//	// AfterField validates if target is after field.
//	func AfterField(target *time.Time, field *time.Time) validate.Field {
//		return validate.Field{
//			Valid:   target != nil && field != nil && target.After(*field),
//			Message: "{0} should be after {1}",
//			Fields:  []any{target, field},
//		}
//	}
//
// Check each package's documentation for more examples and the functions provided.
//
// # Returning responses
//
// The [lit/render] package contains several constructors for implementations of [lit.Response], including JSON
// responses, redirections, no content responses, files and streams.
//
// Custom responses can be created by implementing [lit.Response] or by calling the [lit.ResponseFunc] function. The
// former is preferred if one intends to return the response in multiple handlers. Check the example below.
//
//	// YAMLResponse is a lit.Response that prints a YAML formatted-body as response. It sets
//	// the Content-Type header to "application/x-yaml".
//	//
//	// If the response contains a body but its marshalling fails, YAMLResponse responds an Internal Server Error
//	// with the error message as plain text.
//	type YAMLResponse struct {
//		StatusCode int
//		Header     http.Header
//		Body       any
//	}
//
//	// YAML responds the request with statusCode and a body marshalled as YAML. Nil body equals empty body.
//	//
//	// If body is a string or an error, YAML marshals render.Message with the body assigned to render.Message.Value.
//	// Otherwise, it marshals the body as is.
//	func YAML(statusCode int, body any) YAMLResponse {
//		switch cast := body.(type) {
//		case string:
//			return YAMLResponse{statusCode, make(http.Header), render.Message{Value: cast}}
//		case error:
//			return YAMLResponse{statusCode, make(http.Header), render.Message{Value: cast.Error()}}
//		default:
//			return YAMLResponse{statusCode, make(http.Header), cast}
//		}
//	}
//
//	func (r YAMLResponse) Write(w http.ResponseWriter) {
//		responseHeader := w.Header()
//		for key := range r.Header {
//			responseHeader.Set(key, r.Header.Get(key))
//		}
//
//		if r.Body == nil {
//			w.WriteHeader(r.StatusCode)
//			return
//		}
//
//		w.Header().Set("Content-Type", "application/x-yaml")
//
//		if err := r.writeBody(w); err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//	}
//
//	func (r YAMLResponse) writeBody(w http.ResponseWriter) error {
//		bodyBytes, err := yaml.Marshal(r.Body)
//		if err != nil {
//			return err
//		}
//
//		w.WriteHeader(r.StatusCode)
//
//		_, err = w.Write(bodyBytes)
//
//		return err
//	}
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
//			return YAML(http.StatusBadRequest, err)
//		}
//
//		return YAML(http.StatusOK, req.A/req.B)
//	}
//
// Check the [lit/render] package for more examples and the functions provided.
package lit
