package bind

import (
	"encoding/json"
	"encoding/xml"
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
// Tags from encoding packages, such as "json" tag, can be used appropriately.
//
// If the Content-Type header is not set nor supported, Body defaults to JSON parsing.
func Body[T any](r *lit.Request) (T, error) {
	var target T

	switch r.Header().Get("Content-Type") {
	case "application/xml", "text/xml":
		return target, xml.NewDecoder(r.Body()).Decode(&target)
	case "application/x-yaml":
		return target, yaml.NewDecoder(r.Body()).Decode(&target)
	default:
		return target, json.NewDecoder(r.Body()).Decode(&target)
	}
}
