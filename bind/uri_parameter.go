package bind

import (
	"fmt"
	"github.com/jvcoutinho/lit"
	"golang.org/x/exp/maps"
	"reflect"
	"time"
)

// URIParameter binds a request's URI parameter into a value of type T. T can be either a
// primitive type or a [time.Time].
//
// If the value can't be bound into T, URIParameter returns BindingError.
//
// If parameter is not registered as one of the request's expected parameters, URIParameter panics.
func URIParameter[T primitiveType | time.Time](r *lit.Request, parameter string) (T, error) {
	parameterValue, ok := r.URIParameters()[parameter]
	if !ok {
		panic(fmt.Sprintf("%s has not been defined as one of the request parameters: %v", parameter,
			maps.Keys(r.URIParameters())))
	}

	var target T

	targetValue := reflect.ValueOf(&target).Elem()

	if err := bind(parameterValue, targetValue); err != nil {
		return target, fmt.Errorf("%s: %w", parameter, err)
	}

	return target, nil
}
