// Package bind contains model binding features to be used along [*lit.Request].
//
// # Binding functions
//
// All the functions of this package are binding functions, that parse the string data from the request and binds the
// result to a Go struct. Supported functions are:
//
//   - [URIParameters] - one or more URL parameters;
//   - [URIParameter] - one URL parameter;
//   - [Body] - body or forms;
//   - [Query] - query parameters or GET forms;
//   - [Header] - one or more header fields;
//   - [HeaderField] - one header field;
//   - [Request] - optimized combination of [URIParameters], [Body], [Query] and [Header]. Preferred for reading
//     from multiple inputs.
//
// If any of these functions fails to parse the request (for example, they couldn't bind a non-numeric string into
// an integer field), they return [*Error], that contains a user-friendly message and can be used in the response as is.
//
// # Validations
//
// The bind package is integrated with the [github.com/jvcoutinho/lit/validate] package. If a struct passed as type
// parameter for [Request], [URIParameters], [Body], [Query] or [Header] implements the
// [github.com/jvcoutinho/lit/validate.Validatable] interface with a pointer receiver, these binding functions validate
// the result and can return a validation error with a user-friendly message.
package bind

import (
	"errors"
	"fmt"
	"net/textproto"
	"reflect"
	"strconv"
	"time"

	"golang.org/x/exp/constraints"
)

type primitiveType interface {
	constraints.Ordered | constraints.Complex | bool
}

const nonStructTypeParameter = "T must be a struct type"

type invalidArrayLengthError struct {
	// Maximum expected length for the array.
	ExpectedLength int
	// Actual mismatched length.
	ActualLength int
}

func (e invalidArrayLengthError) Error() string {
	return fmt.Sprintf("expected at most %d elements. Got %d", e.ExpectedLength, e.ActualLength)
}

// Error is returned when a binding to a target value fails.
type Error struct {
	// Incoming value.
	Value string
	// Target of the binding.
	Target reflect.Type
	// The actual error.
	Err error
}

func (e Error) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("%s is not a valid %s", e.Value, e.Target)
	}

	return fmt.Sprintf("%s is not a valid %s: %s", e.Value, e.Target, e.Err)
}

func bind(value string, target reflect.Value) error {
	switch target.Kind() {
	case reflect.String:
		target.SetString(value)
		return nil
	case reflect.Pointer:
		return bindPointer(value, target, bind)
	case reflect.Uint:
		return bindUint(0, value, target)
	case reflect.Uint8:
		return bindUint(8, value, target)
	case reflect.Uint16:
		return bindUint(16, value, target)
	case reflect.Uint32:
		return bindUint(32, value, target)
	case reflect.Uint64:
		return bindUint(64, value, target)
	case reflect.Int:
		return bindInt(0, value, target)
	case reflect.Int8:
		return bindInt(8, value, target)
	case reflect.Int16:
		return bindInt(16, value, target)
	case reflect.Int32:
		return bindInt(32, value, target)
	case reflect.Int64:
		return bindInt(64, value, target)
	case reflect.Float32:
		return bindFloat(32, value, target)
	case reflect.Float64:
		return bindFloat(64, value, target)
	case reflect.Complex64:
		return bindComplex(64, value, target)
	case reflect.Complex128:
		return bindComplex(128, value, target)
	case reflect.Bool:
		return bindBool(value, target)
	case reflect.Struct:
		switch target.Interface().(type) {
		case time.Time:
			return bindTime(value, target)
		}
		fallthrough
	default:
		panic(fmt.Sprintf("unbindable type %s", target.Type()))
	}
}

func bindAll(values []string, target reflect.Value) error {
	switch target.Kind() {
	case reflect.Pointer:
		return bindPointer(values, target, bindAll)
	case reflect.Slice:
		return bindSlice(values, target)
	case reflect.Array:
		return bindArray(values, target)
	default:
		if len(values) == 1 {
			return bind(values[0], target)
		}

		return Error{fmt.Sprint(values), target.Type(), nil}
	}
}

func bindStruct[T any, V string | []string](
	values map[string]V,
	tag string,
	bindFunction func(V, reflect.Value) error,
) (T, error) {
	var target T

	targetValue := reflect.ValueOf(&target).Elem()

	if targetValue.Kind() != reflect.Struct {
		panic(nonStructTypeParameter)
	}

	fields := reflect.VisibleFields(targetValue.Type())

	if err := bindFields(values, tag, targetValue, fields, bindFunction); err != nil {
		return target, err
	}

	if err := validateFields(&target); err != nil {
		return target, err
	}

	return target, nil
}

func bindFields[T string | []string](
	values map[string]T,
	fieldTag string,
	structValue reflect.Value,
	fields []reflect.StructField,
	bindFunction func(T, reflect.Value) error,
) error {
	for _, field := range fields {
		fieldValue := structValue.FieldByIndex(field.Index)

		if !fieldValue.CanSet() {
			continue
		}

		parameter, ok := field.Tag.Lookup(fieldTag)

		if !ok {
			continue
		}

		var value T

		if fieldTag == headerTag {
			value, ok = values[textproto.CanonicalMIMEHeaderKey(parameter)]
		} else {
			value, ok = values[parameter]
		}

		if !ok {
			continue
		}

		if err := bindFunction(value, fieldValue); err != nil {
			return fmt.Errorf("%s: %w", parameter, err)
		}
	}

	return nil
}

func bindPointer[T string | []string](value T, target reflect.Value, fn func(T, reflect.Value) error) error {
	targetValue := reflect.New(target.Type().Elem())
	if err := fn(value, targetValue.Elem()); err != nil {
		return err
	}

	target.Set(targetValue)

	return nil
}

func bindUint(bitSize int, value string, target reflect.Value) error {
	converted, err := strconv.ParseUint(value, 10, bitSize)
	if err != nil {
		return Error{value, target.Type(), errors.Unwrap(err)}
	}

	target.SetUint(converted)

	return nil
}

func bindInt(bitSize int, value string, target reflect.Value) error {
	converted, err := strconv.ParseInt(value, 10, bitSize)
	if err != nil {
		return Error{value, target.Type(), errors.Unwrap(err)}
	}

	target.SetInt(converted)

	return nil
}

func bindFloat(bitSize int, value string, target reflect.Value) error {
	converted, err := strconv.ParseFloat(value, bitSize)
	if err != nil {
		return Error{value, target.Type(), errors.Unwrap(err)}
	}

	target.SetFloat(converted)

	return nil
}

func bindComplex(bitSize int, value string, target reflect.Value) error {
	converted, err := strconv.ParseComplex(value, bitSize)
	if err != nil {
		return Error{value, target.Type(), errors.Unwrap(err)}
	}

	target.SetComplex(converted)

	return nil
}

func bindBool(value string, target reflect.Value) error {
	converted, err := strconv.ParseBool(value)
	if err != nil {
		return Error{value, target.Type(), errors.Unwrap(err)}
	}

	target.SetBool(converted)

	return nil
}

func bindTime(value string, target reflect.Value) error {
	converted, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return Error{value, target.Type(), err}
	}

	target.Set(reflect.ValueOf(converted))

	return nil
}

func bindArray(values []string, target reflect.Value) error {
	if target.Len() < len(values) {
		return Error{fmt.Sprint(values), target.Type(),
			invalidArrayLengthError{target.Len(), len(values)}}
	}

	for i, value := range values {
		if err := bind(value, target.Index(i)); err != nil {
			return err
		}
	}

	return nil
}

func bindSlice(values []string, target reflect.Value) error {
	slice := reflect.MakeSlice(target.Type(), len(values), len(values))
	err := bindArray(values, slice)
	if err != nil {
		return err
	}

	target.Set(slice)

	return nil
}
