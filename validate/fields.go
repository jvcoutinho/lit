package validate

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"unsafe"
)

type notFieldPointerError struct {
	structValue reflect.Value
	validation  int
	target      int
}

func (e notFieldPointerError) Error() string {
	return fmt.Sprintf("validation %d, target %d should be a pointer to a field of %s", e.validation, e.target,
		e.structValue.Type().String())
}

// Fields validate the fields of a struct of type T.
//
// It uses the fields' tag values to display a message for the user in case the validation fails (by
// using binding functions along Validatable, this can be inferred automatically).
// Then, if tag is empty or is missing in the field, it uses the field's name instead.
//
// If T is not a struct type, Fields panics.
func Fields[T any](target *T, tag string, validations ...Field) error {
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
		for ai, argument := range validation.Targets {
			argumentAddress, ok := getArgumentAddress(argumentAddresses, argument)
			if !ok {
				panic(notFieldPointerError{structValue, vi, ai})
			}

			field, ok := fieldsPerAddress[argumentAddress]
			if !ok {
				panic(notFieldPointerError{structValue, vi, ai})
			}

			violations[vi].Message = strings.ReplaceAll(validation.Message,
				fmt.Sprintf("{%d}", ai), getReplacement(field, tag))
		}
	}

	return Error{Violations: violations}
}

func getReplacement(field reflect.StructField, tag string) string {
	if tag == "" {
		return field.Name
	}

	value, ok := field.Tag.Lookup(tag)
	if !ok {
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

		fieldsPerAddress[fieldValue.Addr().UnsafePointer()] = field
	}

	return fieldsPerAddress
}
