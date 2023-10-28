package bind

import (
	"github.com/jvcoutinho/lit"
	"net/http"
	"reflect"
)

// Header binds the request's header into the fields of a struct of type T.
// Targeted fields should be annotated with the tag "header".
//
// If any field couldn't be bound, Header returns BindingError.
//
// If T is not a struct type, Header panics.
func Header[T any](r *lit.Request) (T, error) {
	var target T

	targetValue := reflect.ValueOf(&target).Elem()

	if targetValue.Kind() != reflect.Struct {
		panic(nonStructTypeParameter)
	}

	if err := bindHeader(r.Header(), targetValue.Type(), targetValue); err != nil {
		return target, err
	}

	return target, nil
}

func bindHeader(header http.Header, structType reflect.Type, structValue reflect.Value) error {
	return bindFields[[]string](header, "header", structType, structValue, bindAll)
}
