package bind

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"

	"github.com/jvcoutinho/lit"
	"gopkg.in/yaml.v3"
)

const (
	jsonTag = "json"
	yamlTag = "yaml"
	xmlTag  = "xml"
	formTag = "form"
)

// Body binds the request's body into a value of type T.
//
// It checks the Content-Type header to select an appropriated parsing method:
//   - "application/json" for JSON parsing
//   - "application/xml" or "text/xml" for XML parsing
//   - "application/x-yaml" for YAML parsing
//   - "application/x-www-form-urlencoded" for form parsing (T must be a struct type in this case.
//     Otherwise, Body panics)
//
// Tags from encoding packages, such as "json", "xml" and "yaml" tags, can be used appropriately. For form parsing,
// use the tag "form".
//
// If the Content-Type header is not set nor supported, Body defaults to JSON parsing.
func Body[T any](r *lit.Request) (T, error) {
	var (
		target T
		err    error
	)

	switch r.Header().Get("Content-Type") {
	case "application/xml", "text/xml":
		err = decodeXML(r.Body(), &target)
	case "application/x-yaml", "text/yaml":
		err = decodeYAML(r.Body(), &target)
	case "application/x-www-form-urlencoded":
		return parseForm[T](r)
	default:
		err = decodeJSON(r.Body(), &target)
	}

	if errors.Is(err, io.EOF) {
		return target, nil
	}

	return target, err
}

func parseForm[T any](r *lit.Request) (T, error) {
	var target T

	form, err := r.Form()
	if err != nil {
		return target, err
	}

	return bindStruct[T](form, formTag, bindAll)
}

func decodeJSON(body io.ReadCloser, target any) error {
	err := json.NewDecoder(body).Decode(target)

	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		return fmt.Errorf("%s: %w", unmarshalTypeError.Field,
			BindingError{unmarshalTypeError.Value, unmarshalTypeError.Type, nil},
		)
	}

	return err
}

func decodeYAML(body io.ReadCloser, target any) error {
	return yaml.NewDecoder(body).Decode(target)
}

func decodeXML(body io.ReadCloser, target any) error {
	return xml.NewDecoder(body).Decode(target)
}
