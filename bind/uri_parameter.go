package bind

import (
	"fmt"
	"github.com/jvcoutinho/lit"
	"reflect"
)

// URIParameter binds a request's URI parameter into a value of a simple type T (either a primitive type or a time.Time).
//
// If the value couldn't be bind into T, URIParameter returns a BindingError wrapped with the specific error.
//
// If parameter is not present in the request's defined parameters, URIParameter panics, since this situation
// is generally an implementation error.
func URIParameter[T simpleType](r *lit.Request, parameter string) (T, error) {
	var target T

	targetValue := reflect.ValueOf(&target).Elem()

	parameterValue, ok := r.URIParameters()[parameter]
	if !ok {
		panic(fmt.Sprintf("%s has not been defined", parameter))
	}

	if err := bind(parameterValue, targetValue); err != nil {
		return target, fmt.Errorf("%s: %w", parameter, err)
	}

	return target, nil
}
