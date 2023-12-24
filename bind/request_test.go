package bind_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/jvcoutinho/lit/bind"

	"github.com/jvcoutinho/lit"
	"github.com/stretchr/testify/require"
)

func TestRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description    string
		uriParameters  map[string]string
		query          url.Values
		header         http.Header
		body           string
		function       func(r *lit.Request) (any, error)
		expectedResult any
		expectedError  string
		shouldPanic    bool
	}{
		{
			description: "WhenTypeParameterIsNotStruct_ShouldPanic",
			function: func(r *lit.Request) (any, error) {
				return bind.Request[int](r)
			},
			expectedError: "T must be a struct type",
			shouldPanic:   true,
		},
		{
			description: "WhenTargetHasUnbindableFieldToBindingMethod_ShouldPanic",
			uriParameters: map[string]string{
				"field": "123",
			},
			function: func(r *lit.Request) (any, error) {
				return bind.Request[unbindableField](r)
			},
			expectedError: "unbindable type lit.Request",
			shouldPanic:   true,
		},
		{
			description: "WhenFieldIsUnexported_ShouldIgnore",
			uriParameters: map[string]string{
				"unexported": "123",
			},
			query: url.Values{
				"unexported": {"123"},
			},
			header: http.Header{
				"unexported": {"123"},
			},
			body: `{"unexported": "123"}`,
			function: func(r *lit.Request) (any, error) {
				return bind.Request[ignorableFields](r)
			},
			expectedResult: ignorableFields{},
		},
		{
			description: "WhenFieldIsMissingFromRequest_ShouldIgnore",
			function: func(r *lit.Request) (any, error) {
				return bind.Request[ignorableFields](r)
			},
			expectedResult: ignorableFields{},
		},
		{
			description: "WhenURIParametersFieldsAreValid_ShouldBindThem",
			uriParameters: map[string]string{
				"string":     "hi",
				"pointer":    "10",
				"uint":       "10",
				"uint8":      "10",
				"uint16":     "10",
				"uint32":     "10",
				"uint64":     "10",
				"int":        "10",
				"int8":       "10",
				"int16":      "10",
				"int32":      "10",
				"int64":      "10",
				"float32":    "10.5",
				"float64":    "10.5",
				"complex64":  "10.5",
				"complex128": "10.5",
				"bool":       "true",
				"time":       "2023-10-22T00:00:00Z",
			},
			function: func(r *lit.Request) (any, error) {
				return bind.Request[bindableFields](r)
			},
			expectedResult: bindableFields{
				String:     "hi",
				Pointer:    pointerOf(10),
				Uint:       10,
				Uint8:      10,
				Uint16:     10,
				Uint32:     10,
				Uint64:     10,
				Int:        10,
				Int8:       10,
				Int16:      10,
				Int32:      10,
				Int64:      10,
				Float32:    10.5,
				Float64:    10.5,
				Complex64:  10.5,
				Complex128: 10.5,
				Bool:       true,
				Time:       time.Date(2023, 10, 22, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			description: "WhenQueryParametersFieldsAreValid_ShouldBindThem",
			query: url.Values{
				"string":     {"hi"},
				"pointer":    {"10"},
				"uint":       {"10"},
				"uint8":      {"10"},
				"uint16":     {"10"},
				"uint32":     {"10"},
				"uint64":     {"10"},
				"int":        {"10"},
				"int8":       {"10"},
				"int16":      {"10"},
				"int32":      {"10"},
				"int64":      {"10"},
				"float32":    {"10.5"},
				"float64":    {"10.5"},
				"complex64":  {"10.5"},
				"complex128": {"10.5"},
				"bool":       {"true"},
				"time":       {"2023-10-22T00:00:00Z"},
				"slice":      {"2", "3"},
				"array":      {"2", "3"},
			},
			function: func(r *lit.Request) (any, error) {
				return bind.Request[bindableFields](r)
			},
			expectedResult: bindableFields{
				String:     "hi",
				Pointer:    pointerOf(10),
				Uint:       10,
				Uint8:      10,
				Uint16:     10,
				Uint32:     10,
				Uint64:     10,
				Int:        10,
				Int8:       10,
				Int16:      10,
				Int32:      10,
				Int64:      10,
				Float32:    10.5,
				Float64:    10.5,
				Complex64:  10.5,
				Complex128: 10.5,
				Bool:       true,
				Time:       time.Date(2023, 10, 22, 0, 0, 0, 0, time.UTC),
				Slice:      []int{2, 3},
				Array:      [2]int{2, 3},
			},
		},
		{
			description: "WhenHeaderFieldsAreValid_ShouldBindThem",
			header: http.Header{
				"string":     {"hi"},
				"pointer":    {"10"},
				"uint":       {"10"},
				"uint8":      {"10"},
				"uint16":     {"10"},
				"uint32":     {"10"},
				"uint64":     {"10"},
				"int":        {"10"},
				"int8":       {"10"},
				"int16":      {"10"},
				"int32":      {"10"},
				"int64":      {"10"},
				"float32":    {"10.5"},
				"float64":    {"10.5"},
				"complex64":  {"10.5"},
				"complex128": {"10.5"},
				"bool":       {"true"},
				"time":       {"2023-10-22T00:00:00Z"},
				"slice":      {"2", "3"},
				"array":      {"2", "3"},
			},
			function: func(r *lit.Request) (any, error) {
				return bind.Request[bindableFields](r)
			},
			expectedResult: bindableFields{
				String:     "hi",
				Pointer:    pointerOf(10),
				Uint:       10,
				Uint8:      10,
				Uint16:     10,
				Uint32:     10,
				Uint64:     10,
				Int:        10,
				Int8:       10,
				Int16:      10,
				Int32:      10,
				Int64:      10,
				Float32:    10.5,
				Float64:    10.5,
				Complex64:  10.5,
				Complex128: 10.5,
				Bool:       true,
				Time:       time.Date(2023, 10, 22, 0, 0, 0, 0, time.UTC),
				Slice:      []int{2, 3},
				Array:      [2]int{2, 3},
			},
		},
		{
			description: "WhenBodyFieldsAreValid_ShouldBindThem",
			body: `
{
    "string": "hi",
	"pointer": 10,
    "uint": 10,
    "uint8": 10,
    "uint16": 10,
    "uint32": 10,
    "uint64": 10,
    "int": 10,
    "int8": 10,
    "int16": 10,
    "int32": 10,
    "int64": 10,
    "float32": 10.5,
    "float64": 10.5,
    "bool": true,
    "time": "2023-10-22T00:00:00Z",
    "slice": [
        2,
        3
    ],
    "array": [
        2,
        3
    ]
}`,
			function: func(r *lit.Request) (any, error) {
				return bind.Request[bindableFields](r)
			},
			expectedResult: bindableFields{
				String:  "hi",
				Pointer: pointerOf(10),
				Uint:    10,
				Uint8:   10,
				Uint16:  10,
				Uint32:  10,
				Uint64:  10,
				Int:     10,
				Int8:    10,
				Int16:   10,
				Int32:   10,
				Int64:   10,
				Float32: 10.5,
				Float64: 10.5,
				Bool:    true,
				Time:    time.Date(2023, 10, 22, 0, 0, 0, 0, time.UTC),
				Slice:   []int{2, 3},
				Array:   [2]int{2, 3},
			},
		},
		{
			description:   "WhenURIParameterFieldIsInvalid_ShouldReturnError",
			uriParameters: map[string]string{"uint": "10a"},
			function: func(r *lit.Request) (any, error) {
				return bind.Request[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "uint: 10a is not a valid uint: invalid syntax",
		},
		{
			description: "WhenQueryFieldIsInvalid_ShouldReturnError",
			query:       url.Values{"uint": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Request[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "uint: 10a is not a valid uint: invalid syntax",
		},
		{
			description: "WhenHeaderIsInvalid_ShouldReturnError",
			header:      http.Header{"uint": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Request[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "uint: 10a is not a valid uint: invalid syntax",
		},
		{
			description: "WhenBodyFieldIsInvalid_ShouldReturnError",
			body:        `{"uint": "10a"}`,
			function: func(r *lit.Request) (any, error) {
				return bind.Request[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "uint: string is not a valid uint",
		},
		{
			description: "WhenContentTypeIsYAML_ShouldParseYAMLBody",
			body:        "uint: 10",
			header: http.Header{
				"Content-Type": {"application/x-yaml"},
			},
			function: func(r *lit.Request) (any, error) {
				return bind.Request[bindableFields](r)
			},
			expectedResult: bindableFields{Uint: 10},
		},
		{
			description: "WhenContentTypeIsXML_ShouldParseXMLBody",
			body: `
<?xml version="1.0" encoding="UTF-8" ?>
<root>
  <Uint>10</Uint>
</root>
			`,
			header: http.Header{
				"Content-Type": {"application/xml"},
			},
			function: func(r *lit.Request) (any, error) {
				return bind.Request[bindableFields](r)
			},
			expectedResult: bindableFields{Uint: 10},
		},
		{
			description: "WhenContentTypeIsForm_AndRequestIsPOST_ShouldParseFormBodyAndQueryParameters",
			query:       url.Values{"uint8": {"10"}},
			body:        "uint=10",
			header: http.Header{
				"Content-Type": {"application/x-www-form-urlencoded"},
			},
			function: func(r *lit.Request) (any, error) {
				return bind.Request[bindableFields](r)
			},
			expectedResult: bindableFields{Uint: 10, Uint8: 10},
		},
		{
			description:   "WhenTypeParameterIsValidatableWithValueReceiver_ShouldNotValidate",
			uriParameters: map[string]string{"string": "string"},
			function: func(r *lit.Request) (any, error) {
				return bind.Request[nonPointerReceiverValidatableFields](r)
			},
			expectedResult: nonPointerReceiverValidatableFields{String: "string"},
		},
		{
			description: "WhenTypeParameterIsValidatableWithPointerReceiver_ShouldValidate",
			body:        `{"string": "string"}`,
			function: func(r *lit.Request) (any, error) {
				return bind.Request[pointerReceiverValidatableFields](r)
			},
			expectedResult: pointerReceiverValidatableFields{String: "string"},
			expectedError:  "string should have a length greater than 6",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(test.body))
			request.URL.RawQuery = test.query.Encode()
			for key, value := range test.header {
				for _, v := range value {
					request.Header.Add(key, v)
				}
			}

			r := lit.NewRequest(request).WithURIParameters(test.uriParameters)

			// Act
			if test.shouldPanic {
				require.PanicsWithValue(t, test.expectedError, func() {
					_, _ = test.function(r)
				})

				return
			}

			result, err := test.function(r)

			// Assert
			errMessage := ""
			if err != nil {
				errMessage = err.Error()
			}

			require.Equal(t, test.expectedError, errMessage)
			require.Equal(t, test.expectedResult, result)
		})
	}
}
