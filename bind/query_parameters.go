package bind

import (
	"fmt"
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
		panic(fmt.Sprintf("%T is not a struct type", target))
	}

	err := bindQueryParameters(r.URL().Query(), targetValue.Type(), targetValue)
	if err != nil {
		return target, err
	}

	return target, nil
}

func bindQueryParameters(parameters url.Values, structType reflect.Type, structValue reflect.Value) error {
	numberFields := structType.NumField()

	for i := 0; i < numberFields; i++ {
		fieldType := structType.Field(i)

		if !fieldType.IsExported() {
			continue
		}

		parameter, ok := fieldType.Tag.Lookup("query")

		if !ok {
			continue
		}

		values, ok := parameters[parameter]

		if !ok {
			continue
		}

		if err := bindAll(values, structValue.Field(i)); err != nil {
			return fmt.Errorf("%s: %w", parameter, err)
		}
	}

	return nil
}
