package bind

import (
	"github.com/jvcoutinho/lit"
)

const uriParameterTag = "uri"

// URIParameters binds the request's URI parameters into the fields of a struct of type T.
// Targeted fields should be exported and annotated with the tag "uri". Otherwise, they are ignored.
//
// If a field can't be bound, URIParameters returns Error.
//
// If *T implements validate.Validatable (with a pointer receiver), URIParameters calls validate.Fields on the result
// and can return validate.Error.
//
// If T is not a struct type, URIParameters panics.
func URIParameters[T any](r *lit.Request) (T, error) {
	return bindStruct[T](r.URIParameters(), uriParameterTag, bind)
}
