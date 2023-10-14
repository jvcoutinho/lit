package bind

import (
	"fmt"
	"reflect"
	"strconv"
)

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
		return bindBoolInto(value, target)
	default:
		panic(fmt.Sprintf("unsupported type %s", target.Kind()))
	}
}

func bindUint(bitSize int, value string, target reflect.Value) error {
	converted, err := strconv.ParseUint(value, 10, bitSize)
	if err != nil {
		if bitSize == 0 {
			return fmt.Errorf("parsing %s into uint: %w", value, err)
		}

		return fmt.Errorf("parsing %s into uint%d: %w", value, bitSize, err)
	}

	target.SetUint(converted)

	return nil
}

func bindInt(bitSize int, value string, target reflect.Value) error {
	converted, err := strconv.ParseInt(value, 10, bitSize)
	if err != nil {
		return fmt.Errorf("parsing %s into %s: %w", value, target.Kind(), err)
	}

	target.SetInt(converted)

	return nil
}

func bindFloat(bitSize int, value string, target reflect.Value) error {
	converted, err := strconv.ParseFloat(value, bitSize)
	if err != nil {
		return fmt.Errorf("parsing %s into float%d: %w", value, bitSize, err)
	}

	target.SetFloat(converted)

	return nil
}

func bindComplex(bitSize int, value string, target reflect.Value) error {
	converted, err := strconv.ParseComplex(value, bitSize)
	if err != nil {
		return fmt.Errorf("parsing %s into complex%d: %w", value, bitSize, err)
	}

	target.SetComplex(converted)

	return nil
}

func bindBoolInto(value string, target reflect.Value) error {
	converted, err := strconv.ParseBool(value)
	if err != nil {
		return fmt.Errorf("parsing %s into bool: %w", value, err)
	}

	target.SetBool(converted)

	return nil
}
