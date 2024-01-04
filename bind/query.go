package bind

import (
	"github.com/jvcoutinho/lit"
)

const queryParameterTag = "query"

// Query binds the request's query parameters into the fields of a struct of type T.
// Targeted fields should be exported and annotated with the tag "query". Otherwise, they are ignored.
//
// If a field can't be bound, Query returns Error.
//
// If *T implements [validate.Validatable] (with a pointer receiver), Query calls [validate.Fields] on the result
// and can return [validate.Error].
//
// If T is not a struct type, Query panics.
func Query[T any](r *lit.Request) (T, error) {
	return bindStruct[T](r.URL().Query(), queryParameterTag, bindAll)
}
