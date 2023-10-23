package bind

import (
	"fmt"
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
		panic(fmt.Sprintf("%T is not a struct type", target))
	}

	err := bindURLParameters(r.URLParameters(), targetValue.Type(), targetValue)
	if err != nil {
		return target, err
	}

	return target, nil
}

func bindURLParameters(parameters map[string]string, structType reflect.Type, structValue reflect.Value) error {
	numberFields := structType.NumField()

	for i := 0; i < numberFields; i++ {
		fieldType := structType.Field(i)

		parameter, ok := fieldType.Tag.Lookup("uri")

		if !ok {
			continue
		}

		argument, ok := parameters[parameter]

		if !ok {
			continue
		}

		if err := bind(argument, structValue.Field(i)); err != nil {
			return fmt.Errorf("%s: %w", parameter, err)
		}
	}

	return nil
}
