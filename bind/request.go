package bind

import (
	"net/http"
	"reflect"

	"github.com/jvcoutinho/lit"
)

// Request binds the request's body, query, header and URI parameters into the fields of a struct of type T.
// Targeted fields should be exported and annotated with corresponding binding tags. Otherwise, they are ignored.
//
// It's an optimized combination of the binding functions [Body], [Query], [Header] and [URIParameters], suitable
// when you need to read from multiple inputs of the request.
//
// If a field can't be bound, Request returns Error.
//
// If *T implements [validate.Validatable] (with a pointer receiver), Request calls [validate.Fields] on the result
// and can return [validate.Error].
//
// If T is not a struct type, Request panics.
func Request[T any](r *lit.Request) (T, error) {
	var target T

	targetValue := reflect.ValueOf(&target).Elem()

	if targetValue.Kind() != reflect.Struct {
		panic(nonStructTypeParameter)
	}

	var (
		fields       = reflect.VisibleFields(targetValue.Type())
		fieldsPerTag = getFieldsPerTag(fields)
	)

	var (
		uriParameters, uriParametersFields = r.URIParameters(), fieldsPerTag[uriParameterTag]
		query, queryFields                 = r.URL().Query(), fieldsPerTag[queryParameterTag]
		header, headerFields               = r.Header(), fieldsPerTag[headerTag]
		body                               = r.Body()
	)

	if body != http.NoBody {
		if err := bindBody(r, &target, targetValue); err != nil {
			return target, err
		}
	}

	if len(uriParameters) > 0 && len(uriParametersFields) > 0 {
		if err := bindFields(r.URIParameters(), uriParameterTag, targetValue, uriParametersFields, bind); err != nil {
			return target, err
		}
	}

	if len(query) > 0 && len(queryFields) > 0 {
		if err := bindFields(r.URL().Query(), queryParameterTag, targetValue, queryFields, bindAll); err != nil {
			return target, err
		}
	}

	if len(header) > 0 && len(headerFields) > 0 {
		if err := bindFields(r.Header(), headerTag, targetValue, headerFields, bindAll); err != nil {
			return target, err
		}
	}

	if err := validateFields(&target); err != nil {
		return target, err
	}

	return target, nil
}

func getFieldsPerTag(fields []reflect.StructField) map[string][]reflect.StructField {
	fieldsPerTag := map[string][]reflect.StructField{}
	tags := []string{uriParameterTag, queryParameterTag, headerTag}

	for _, field := range fields {
		for _, tag := range tags {
			if _, ok := field.Tag.Lookup(tag); ok {
				fieldsPerTag[tag] = append(fieldsPerTag[tag], field)
			}
		}
	}

	return fieldsPerTag
}
