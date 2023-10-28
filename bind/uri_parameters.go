package bind

import (
	"github.com/jvcoutinho/lit"
	"reflect"
)

// URIParameters binds the request's URI parameters into the values of a struct of type T.
// Targeted fields should be annotated with the tag "uri".
//
// If T is not a struct type, URIParameters panics.
func URIParameters[T any](r *lit.Request) (T, error) {
	var target T

	targetValue := reflect.ValueOf(&target).Elem()

	if targetValue.Kind() != reflect.Struct {
		panic(nonStructTypeParameter)
	}

	err := bindURIParameters(r.URIParameters(), targetValue.Type(), targetValue)
	if err != nil {
		return target, err
	}

	return target, nil
}

func bindURIParameters(parameters map[string]string, structType reflect.Type, structValue reflect.Value) error {
	return bindFields[string](parameters, "uri", structType, structValue, bind)
}
