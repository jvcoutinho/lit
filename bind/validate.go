package bind

import (
	"log"
	"reflect"
	"unsafe"

	"github.com/jvcoutinho/lit/validate"
)

func validateFields[T any](target T, tag string) error {
	validatable, ok := any(&target).(validate.Validatable)
	if !ok { // we want pointer receiver to implement the interface
		return nil
	}

	_, ok = any(target).(validate.Validatable)
	if ok { // but we don't want value receiver to implement the interface
		log.Printf("%T: the receiver of Validate() should be a pointer in order to use validation features", target)
		return nil
	}

	validation := validatable.Validate()

	if validation.Valid {
		return nil
	}

	structValue := reflect.ValueOf(validatable).Elem()

	fields := reflect.VisibleFields(structValue.Type())

	fieldsPerAddress := getFieldsPerAddress(fields, structValue)

	replacements := make([]any, len(validation.Arguments))
	for i, argument := range validation.Arguments {
		argumentValue := reflect.ValueOf(argument)

		if argumentValue.Kind() == reflect.Pointer {
			argumentAddress := argumentValue.UnsafePointer()

			if field, ok := fieldsPerAddress[argumentAddress]; ok {
				replacements[i] = field.Tag.Get(tag)

				continue
			}
		}

		replacements[i] = argument
	}

	return validate.Validation{
		Valid:     validation.Valid,
		Format:    validation.Format,
		Arguments: replacements,
	}
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
