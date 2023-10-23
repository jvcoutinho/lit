// Package bind contains model binding features to be used along *lit.Request.
package bind

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// ErrUnsupportedType is returned when a bind to an unsupported type is attempted.
var ErrUnsupportedType = errors.New("unsupported type for binding")

// BindingError occurs when a binding is not possible due to type compatibility.
// For example, by trying to bind an integer into a boolean variable.
type BindingError struct {
	// Incoming value.
	Value string
	// Target of the binding.
	Target reflect.Type
}

func (e BindingError) Error() string {
	return fmt.Sprintf("%s is not a valid %s", e.Value, e.Target)
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
		return ErrUnsupportedType
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

		return ErrUnsupportedType
	}
}

func bindUint(bitSize int, value string, target reflect.Value) error {
	converted, err := strconv.ParseUint(value, 10, bitSize)
	if err != nil {
		return fmt.Errorf("%s: %w", BindingError{value, target.Type()}, errors.Unwrap(err))
	}

	target.SetUint(converted)

	return nil
}

func bindInt(bitSize int, value string, target reflect.Value) error {
	converted, err := strconv.ParseInt(value, 10, bitSize)
	if err != nil {
		return fmt.Errorf("%s: %w", BindingError{value, target.Type()}, errors.Unwrap(err))
	}

	target.SetInt(converted)

	return nil
}

func bindFloat(bitSize int, value string, target reflect.Value) error {
	converted, err := strconv.ParseFloat(value, bitSize)
	if err != nil {
		return fmt.Errorf("%s: %w", BindingError{value, target.Type()}, errors.Unwrap(err))
	}

	target.SetFloat(converted)

	return nil
}

func bindComplex(bitSize int, value string, target reflect.Value) error {
	converted, err := strconv.ParseComplex(value, bitSize)
	if err != nil {
		return fmt.Errorf("%s: %w", BindingError{value, target.Type()}, errors.Unwrap(err))
	}

	target.SetComplex(converted)

	return nil
}

func bindBool(value string, target reflect.Value) error {
	converted, err := strconv.ParseBool(value)
	if err != nil {
		return fmt.Errorf("%s: %w", BindingError{value, target.Type()}, errors.Unwrap(err))
	}

	target.SetBool(converted)

	return nil
}

func bindTime(value string, target reflect.Value) error {
	converted, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return fmt.Errorf("%s: %w", BindingError{value, target.Type()}, err)
	}

	target.Set(reflect.ValueOf(converted))

	return nil
}

func bindArray(values []string, target reflect.Value) error {
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
