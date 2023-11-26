package validate

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"unsafe"
)

const validateTag = "validate"

type notFieldPointerError struct {
	structValue reflect.Value
	target      int
}

func (e notFieldPointerError) Error() string {
	return fmt.Sprintf("argument %d should be a pointer to a field of %s", e.target,
		e.structValue.Type().String())
}

// Fields validate the fields of a struct of type T.
//
// It uses the "validate" tag from the fields to build a message for the user in case the validation fails.
// Then, if the tag is the empty string or is missing in a field, it uses the field's name instead.
//
// If T is not a struct type, Fields panics.
func Fields[T any](target *T, validations ...Field) error {
	violations := slices.DeleteFunc(validations, func(f Field) bool { return f.Valid })

	if len(violations) == 0 {
		return nil
	}

	var (
		structValue       = reflect.ValueOf(target).Elem()
		fields            = reflect.VisibleFields(structValue.Type())
		fieldsPerAddress  = getFieldsPerAddress(fields, structValue)
		argumentAddresses = make(map[any]unsafe.Pointer)
	)

	for vi, validation := range violations {
		for ai, argument := range validation.Fields {
			argumentAddress, ok := getArgumentAddress(argumentAddresses, argument)
			if !ok {
				panic(notFieldPointerError{structValue, ai})
			}

			field, ok := fieldsPerAddress[argumentAddress]
			if !ok {
				panic(notFieldPointerError{structValue, ai})
			}

			violations[vi].Message = strings.ReplaceAll(validation.Message,
				fmt.Sprintf("{%d}", ai), getReplacement(field, validateTag))
		}
	}

	return Error{Violations: violations}
}

func getReplacement(field reflect.StructField, tag string) string {
	value, ok := field.Tag.Lookup(tag)
	if !ok || value == "" {
		return field.Name
	}

	return value
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
