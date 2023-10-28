package bind

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/jvcoutinho/lit"
	"gopkg.in/yaml.v3"
)

// Body binds the request's body into a value of type T.
//
// It checks the Content-Type header to select an appropriated parsing method:
//   - "application/json" for JSON parsing
//   - "application/xml" or "text/xml" for XML parsing
//   - "application/x-yaml" for YAML parsing
//
// Tags from encoding packages, such as "json", "xml" and "yaml" tags, can be used appropriately.
//
// If the Content-Type header is not set nor supported, Body defaults to JSON parsing.
func Body[T any](r *lit.Request) (T, error) {
	var target T

	switch r.Header().Get("Content-Type") {
	case "application/xml", "text/xml":
		if err := xml.NewDecoder(r.Body()).Decode(&target); err != nil {
			return target, fmt.Errorf("invalid XML: %w", err)
		}
	case "application/x-yaml", "text/yaml":
		if err := yaml.NewDecoder(r.Body()).Decode(&target); err != nil {
			return target, fmt.Errorf("invalid YAML: %w", err)
		}
	default:
		if err := json.NewDecoder(r.Body()).Decode(&target); err != nil {
			return target, fmt.Errorf("invalid JSON: %w", err)
		}
	}

	return target, nil
}
