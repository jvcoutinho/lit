package bind

import (
	"net/http"
	"reflect"

	"github.com/jvcoutinho/lit"
)

// Request binds the request's body, query, header and URI parameters into the fields of a struct of type T.
// Targeted fields should be exported and annotated with corresponding binding tags.
//
// It's an optimized combination of the binding functions [Body], [Query], [Header] and [URIParameters], suitable
// when you need to read from multiple inputs of the request.
//
// If a field can't be bound, Request returns BindingError.
//
// If T is not a struct type, Request panics.
func Request[T any](r *lit.Request) (T, error) {
	var (
		target T
		err    error
	)

	targetValue := reflect.ValueOf(&target).Elem()

	if targetValue.Kind() != reflect.Struct {
		panic(nonStructTypeParameter)
	}

	fieldsPerTag := getFieldsPerTag(targetValue)

	var (
		uriParameters, uriParametersFields = r.URIParameters(), fieldsPerTag[uriParameterTag]
		query, queryFields                 = r.URL().Query(), fieldsPerTag[queryParameterTag]
		header, headerFields               = r.Header(), fieldsPerTag[headerTag]
		body, hasBodyFields                = r.Body(),
			len(fieldsPerTag[jsonTag]) > 0 || len(fieldsPerTag[yamlTag]) > 0 || len(fieldsPerTag[xmlTag]) > 0
	)

	if body != http.NoBody && hasBodyFields {
		if target, err = Body[T](r); err != nil {
			return target, err
		}
	}

	if len(uriParameters) > 0 && len(uriParametersFields) > 0 {
		if err = bindFields(r.URIParameters(), uriParameterTag, targetValue, uriParametersFields, bind); err != nil {
			return target, err
		}
	}

	if len(query) > 0 && len(queryFields) > 0 {
		if err = bindFields(r.URL().Query(), queryParameterTag, targetValue, queryFields, bindAll); err != nil {
			return target, err
		}
	}

	if len(header) > 0 && len(headerFields) > 0 {
		if err = bindFields(r.Header(), headerTag, targetValue, headerFields, bindAll); err != nil {
			return target, err
		}
	}

	return target, nil
}

func getFieldsPerTag(structValue reflect.Value) map[string][]reflect.StructField {
	fieldsPerTag := map[string][]reflect.StructField{}

	fields := reflect.VisibleFields(structValue.Type())

	for _, field := range fields {
		appendIfContainsTags(fieldsPerTag, field,
			uriParameterTag, queryParameterTag, headerTag, jsonTag, yamlTag, xmlTag)
	}

	return fieldsPerTag
}

func appendIfContainsTags(fieldsPerTag map[string][]reflect.StructField, field reflect.StructField, tags ...string) {
	for _, tag := range tags {
		if _, ok := field.Tag.Lookup(tag); ok {
			fieldsPerTag[tag] = append(fieldsPerTag[tag], field)
		}
	}
}
