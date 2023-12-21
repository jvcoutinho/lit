package validate

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"unsafe"
)

const (
	validateTag            = "validate"
	jsonTag                = "json"
	uriTag                 = "uri"
	queryTag               = "query"
	headerTag              = "header"
	formTag                = "form"
	nonStructTypeParameter = "T must be a struct type"
)

type notFieldPointerError struct {
	structValue reflect.Value
	target      int
}

func (e notFieldPointerError) Error() string {
	return fmt.Sprintf("argument %d should be a pointer to a field of %s", e.target,
		e.structValue.Type().String())
}

// Fields validates the fields of a struct of type T.
//
// It uses the "validate" tag from the fields to build a message for the user in case the validation fails.
// If the tag is set as the empty string or is missing in a field, it tries to use the value from the tags of the
// binding functions.
// If none are present, it uses the field's name instead.
//
// If T is not a struct type, Fields panics.
func Fields[T any](target *T, validations ...Field) error {
	violations := slices.DeleteFunc(validations, func(f Field) bool { return f.Valid })

	if len(violations) == 0 {
		return nil
	}

	structValue := reflect.ValueOf(target).Elem()

	if structValue.Kind() != reflect.Struct {
		panic(nonStructTypeParameter)
	}

	var (
		fields            = reflect.VisibleFields(structValue.Type())
		fieldsPerAddress  = getFieldsPerAddress(fields, structValue)
		argumentAddresses = make(map[any]unsafe.Pointer)
	)

	for vi, validation := range violations {
		for ai, argument := range validation.Fields {
			argumentAddress, ok := getArgumentAddress(argumentAddresses, argument)
			if !ok {
				panic(notFieldPointerError{structValue, ai}.Error())
			}

			field, ok := fieldsPerAddress[argumentAddress]
			if !ok {
				panic(notFieldPointerError{structValue, ai}.Error())
			}

			violations[vi].Message = strings.ReplaceAll(
				violations[vi].Message,
				fmt.Sprintf("{%d}", ai),
				getReplacement(field, validateTag, jsonTag, uriTag, queryTag, headerTag, formTag),
			)
		}
	}

	return Error{Violations: violations}
}

func getReplacement(field reflect.StructField, tags ...string) string {
	for _, tag := range tags {
		value, ok := field.Tag.Lookup(tag)
		if ok && value != "" {
			return value
		}
	}

	return field.Name
}

func getArgumentAddress(addresses map[any]unsafe.Pointer, argument any) (unsafe.Pointer, bool) {
	address, ok := addresses[argument]
	if ok {
		return address, true
	}

	argumentValue := reflect.ValueOf(argument)

	if argumentValue.Kind() != reflect.Pointer {
		return nil, false
	}

	address, ok = argumentValue.UnsafePointer(), true

	addresses[argument] = address

	return address, ok
}

func getFieldsPerAddress(
	fields []reflect.StructField,
	structValue reflect.Value,
) map[unsafe.Pointer]reflect.StructField {
	fieldsPerAddress := make(map[unsafe.Pointer]reflect.StructField, len(fields))

	for _, field := range fields {
		fieldValue := structValue.FieldByIndex(field.Index)

		if fieldValue.Kind() == reflect.Pointer {
			fieldsPerAddress[fieldValue.UnsafePointer()] = field
		} else {
			fieldsPerAddress[fieldValue.Addr().UnsafePointer()] = field
		}
	}

	return fieldsPerAddress
}
