// Package bind contains model binding features to be used along *lit.Request.
package bind

import (
	"errors"
	"fmt"
	"golang.org/x/exp/constraints"
	"reflect"
	"strconv"
	"time"
)

type simpleType interface {
	constraints.Ordered | constraints.Complex | bool | time.Time
}

const nonStructTypeParameter = "T must be a struct type"

// InvalidArrayLengthError is returned when a binding to an array value fails due to mismatched length.
type InvalidArrayLengthError struct {
	// Maximum expected length for the array.
	ExpectedLength int
	// Actual mismatched length.
	ActualLength int
}

func (e InvalidArrayLengthError) Error() string {
	return fmt.Sprintf("expected at most %d elements. Got %d", e.ExpectedLength, e.ActualLength)
}

// BindingError is returned when a binding to a target value fails.
type BindingError struct {
	// Incoming value.
	Value string
	// Target of the binding.
	Target reflect.Type
	// The actual error.
	Err error
}

func (e BindingError) Error() string {
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
		panic(fmt.Sprintf("unsupported type %s", target.Type().Name()))
	}
}

func bindAll(values []string, target reflect.Value) error {
	switch target.Kind() {
	case reflect.Slice:
		return bindSlice(values, target)
	case reflect.Array:
		return bindArray(values, target)
	default:
		if len(values) == 1 {
			return bind(values[0], target)
		}

		return BindingError{fmt.Sprint(values), target.Type(), nil}
	}
}

func bindFields[T string | []string](
	values map[string]T,
	fieldTag string,
	structType reflect.Type,
	structValue reflect.Value,
	bindFunction func(T, reflect.Value) error,
) error {
	numberFields := structType.NumField()

	for i := 0; i < numberFields; i++ {
		fieldType := structType.Field(i)

		if !fieldType.IsExported() {
			continue
		}

		parameter, ok := fieldType.Tag.Lookup(fieldTag)

		if !ok {
			continue
		}

		value, ok := values[parameter]

		if !ok {
			continue
		}

		if err := bindFunction(value, structValue.Field(i)); err != nil {
			return fmt.Errorf("%s: %w", parameter, err)
		}
	}

	return nil
}

func bindUint(bitSize int, value string, target reflect.Value) error {
	converted, err := strconv.ParseUint(value, 10, bitSize)
	if err != nil {
		return BindingError{value, target.Type(), errors.Unwrap(err)}
	}

	target.SetUint(converted)

	return nil
}

func bindInt(bitSize int, value string, target reflect.Value) error {
	converted, err := strconv.ParseInt(value, 10, bitSize)
	if err != nil {
		return BindingError{value, target.Type(), errors.Unwrap(err)}
	}

	target.SetInt(converted)

	return nil
}

func bindFloat(bitSize int, value string, target reflect.Value) error {
	converted, err := strconv.ParseFloat(value, bitSize)
	if err != nil {
		return BindingError{value, target.Type(), errors.Unwrap(err)}
	}

	target.SetFloat(converted)

	return nil
}

func bindComplex(bitSize int, value string, target reflect.Value) error {
	converted, err := strconv.ParseComplex(value, bitSize)
	if err != nil {
		return BindingError{value, target.Type(), errors.Unwrap(err)}
	}

	target.SetComplex(converted)

	return nil
}

func bindBool(value string, target reflect.Value) error {
	converted, err := strconv.ParseBool(value)
	if err != nil {
		return BindingError{value, target.Type(), errors.Unwrap(err)}
	}

	target.SetBool(converted)

	return nil
}

func bindTime(value string, target reflect.Value) error {
	converted, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return BindingError{value, target.Type(), err}
	}

	target.Set(reflect.ValueOf(converted))

	return nil
}

func bindArray(values []string, target reflect.Value) error {
	if target.Len() < len(values) {
		return BindingError{fmt.Sprint(values), target.Type(),
			InvalidArrayLengthError{target.Len(), len(values)}}
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
