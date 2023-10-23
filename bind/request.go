package bind

import (
	"fmt"
	"github.com/jvcoutinho/lit"
	"reflect"
)

// Request binds the request's body, query parameters, header and URL parameters into the fields of a struct of type T.
//
// It is optimized to not do unnecessary computations, so it is appropriate for general usage.
//
// If T is not a struct type, Request panics.
func Request[T any](r *lit.Request) (T, error) {
	target, err := Body[T](r)
	if err != nil {
		return target, err
	}

	targetValue := reflect.ValueOf(&target).Elem()

	if targetValue.Kind() != reflect.Struct {
		panic(fmt.Sprintf("%T is not a struct type", target))
	}

	targetType := targetValue.Type()

	urlParameters := r.URLParameters()
	if len(urlParameters) > 0 {
		if err = bindURLParameters(urlParameters, targetType, targetValue); err != nil {
			return target, err
		}
	}

	queryParameters := r.URL().Query()
	if len(queryParameters) > 0 {
		if err = bindQueryParameters(queryParameters, targetType, targetValue); err != nil {
			return target, err
		}
	}

	header := r.Header()
	if len(header) > 0 {
		if err = bindHeader(header, targetType, targetValue); err != nil {
			return target, err
		}
	}

	return target, nil
}
