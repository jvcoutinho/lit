package bind

import (
	"fmt"
	"github.com/jvcoutinho/lit/validate"
	"log"
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

func validateFields[T any](target *T, targetValue reflect.Value, fields []reflect.StructField, tag string) error {
	validatable, ok := any(target).(validate.Validatable)
	if !ok { // we want pointer receiver to implement the interface
		return nil
	}

	_, ok = any(*target).(validate.Validatable)
	if ok { // but we don't want value receiver to implement the interface
		log.Printf("%T: the receiver of Validate() should be a pointer in order to use "+
			"validation from bind functions", *target)
		return nil
	}

	var (
		validations = validatable.Validate()
		violations  = slices.DeleteFunc(validations, func(f validate.Field) bool { return f.Valid })
	)

	if len(violations) == 0 {
		return nil
	}

	var (
		fieldsPerAddress  = getFieldsPerAddress(fields, targetValue)
		argumentAddresses = make(map[any]unsafe.Pointer)
	)

	for vi, validation := range violations {
		for ai, argument := range validation.Targets {
			argumentAddress, ok := getArgumentAddress(argumentAddresses, argument)
			if !ok {
				panic(notFieldPointerError{targetValue, vi, ai})
			}

			field, ok := fieldsPerAddress[argumentAddress]
			if !ok {
				panic(notFieldPointerError{targetValue, vi, ai})
			}

			violations[vi].Message = strings.ReplaceAll(validation.Message,
				fmt.Sprintf("{%d}", ai), getReplacement(field, tag))
		}
	}

	return validate.Error{Violations: violations}
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
