package bind

import (
	"github.com/jvcoutinho/lit"
	"net/url"
	"reflect"
)

// Query binds the request's query parameters into the values of a struct of type T.
// Targeted fields should be annotated with the tag "query".
//
// If T is not a struct type, Query panics.
func Query[T any](r *lit.Request) (T, error) {
	var target T

	targetValue := reflect.ValueOf(&target).Elem()

	if targetValue.Kind() != reflect.Struct {
		panic(nonStructTypeParameter)
	}

	err := bindQueryParameters(r.URL().Query(), targetValue.Type(), targetValue)
	if err != nil {
		return target, err
	}

	return target, nil
}

func bindQueryParameters(parameters url.Values, structType reflect.Type, structValue reflect.Value) error {
	return bindFields[[]string](parameters, "query", structType, structValue, bindAll)
}
