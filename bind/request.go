package bind

import (
	"errors"
	"github.com/jvcoutinho/lit"
	"io"
	"net/http"
	"reflect"
)

// Request binds the request's body, query parameters, header and URL parameters into the fields of a struct of type T.
//
// If T is not a struct type, Request panics.
func Request[T any](r *lit.Request) (T, error) {
	var (
		target T
		err    error
	)

	targetValue := reflect.ValueOf(&target).Elem()

	if targetValue.Kind() != reflect.Struct {
		panic("T must be a struct type")
	}

	if r.Body() != http.NoBody {
		target, err = Body[T](r)
		if err != nil && !errors.Is(err, io.EOF) {
			return target, err
		}
	}

	targetType := targetValue.Type()

	urlParameters := r.URLParameters()
	if len(urlParameters) > 0 {
		if err := bindURLParameters(urlParameters, targetType, targetValue); err != nil {
			return target, err
		}
	}

	queryParameters := r.URL().Query()
	if len(queryParameters) > 0 {
		if err := bindQueryParameters(queryParameters, targetType, targetValue); err != nil {
			return target, err
		}
	}

	header := r.Header()
	if len(header) > 0 {
		if err := bindHeader(header, targetType, targetValue); err != nil {
			return target, err
		}
	}

	return target, nil
}
