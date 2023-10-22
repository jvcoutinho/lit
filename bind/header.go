package bind

import (
	"fmt"
	"github.com/jvcoutinho/lit"
	"net/http"
	"reflect"
)

// Header binds the request's header into the fields of a struct of type T.
// Targeted fields should be annotated with the tag "header".
//
// If T is not a struct type, Header panics.
func Header[T any](r *lit.Request) (T, error) {
	var target T

	targetValue := reflect.ValueOf(&target).Elem()

	if targetValue.Kind() != reflect.Struct {
		panic(fmt.Sprintf("%T is not a struct type", target))
	}

	err := bindHeader(r.Header(), targetValue.Type(), targetValue)
	if err != nil {
		return target, err
	}

	return target, nil
}

func bindHeader(header http.Header, structType reflect.Type, structValue reflect.Value) error {
	numberFields := structType.NumField()

	for i := 0; i < numberFields; i++ {
		fieldType := structType.Field(i)

		parameter, ok := fieldType.Tag.Lookup("header")

		if !ok {
			continue
		}

		values, ok := header[parameter]

		if !ok {
			continue
		}

		if err := bindAll(values, structValue.Field(i)); err != nil {
			return fmt.Errorf("%s: %w", parameter, err)
		}
	}

	return nil
}
