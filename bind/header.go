package bind

import (
	"github.com/jvcoutinho/lit"
)

const headerTag = "header"

// Header binds the request's header into the fields of a struct of type T.
// Targeted fields should be exported and annotated with the tag "header" (case-insensitive). Otherwise, they are
// ignored.
//
// If any field couldn't be bound, Header returns Error.
//
// If T is not a struct type, Header panics.
func Header[T any](r *lit.Request) (T, error) {
	return bindStruct[T](r.Header(), headerTag, bindAll)
}
