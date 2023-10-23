package bind

import (
	"github.com/jvcoutinho/lit"
	"reflect"
)

// URLParameters binds the request's URL arguments into the values of a struct of type T.
// Targeted fields should be annotated with the tag "uri".
//
// If T is not a struct type, URLParameters panics.
func URLParameters[T any](r *lit.Request) (T, error) {
	var target T

	targetValue := reflect.ValueOf(&target).Elem()

	if targetValue.Kind() != reflect.Struct {
		panic(nonStructTypeParameter)
	}

	err := bindURLParameters(r.URLParameters(), targetValue.Type(), targetValue)
	if err != nil {
		return target, err
	}

	return target, nil
}

func bindURLParameters(parameters map[string]string, structType reflect.Type, structValue reflect.Value) error {
	return bindFields[string](parameters, "uri", structType, structValue, bind)
}
