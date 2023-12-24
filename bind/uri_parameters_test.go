package bind_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/bind"
	"github.com/stretchr/testify/require"
)

func TestURIParameters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description    string
		parameters     map[string]string
		function       func(r *lit.Request) (any, error)
		expectedResult any
		expectedError  string
		shouldPanic    bool
	}{
		{
			description: "WhenTypeParameterIsNotStruct_ShouldPanic",
			function: func(r *lit.Request) (any, error) {
				return bind.URIParameters[int](r)
			},
			expectedError: "T must be a struct type",
			shouldPanic:   true,
		},
		{
			description: "WhenTargetHasUnbindableField_ShouldPanic",
			parameters: map[string]string{
				"field": "123",
			},
			function: func(r *lit.Request) (any, error) {
				return bind.URIParameters[unbindableField](r)
			},
			expectedError: "unbindable type lit.Request",
			shouldPanic:   true,
		},
		{
			description: "WhenFieldIsUnexported_ShouldIgnore",
			parameters: map[string]string{
				"unexported": "123",
			},
			function: func(r *lit.Request) (any, error) {
				return bind.URIParameters[ignorableFields](r)
			},
			expectedResult: ignorableFields{},
		},
		{
			description: "WhenFieldIsMissingFromRequest_ShouldIgnore",
			parameters:  map[string]string{},
			function: func(r *lit.Request) (any, error) {
				return bind.URIParameters[ignorableFields](r)
			},
			expectedResult: ignorableFields{},
		},
		{
			description: "WhenFieldsAreValid_ShouldBindThem",
			parameters: map[string]string{
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
				return bind.URIParameters[bindableFields](r)
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
			description: "WhenPointerValueIsInvalid_ShouldReturnError",
			parameters:  map[string]string{"pointer": "10a"},
			function: func(r *lit.Request) (any, error) {
				return bind.URIParameters[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "pointer: 10a is not a valid int: invalid syntax",
		},
		{
			description: "WhenUintIsInvalid_ShouldReturnError",
			parameters:  map[string]string{"uint": "10a"},
			function: func(r *lit.Request) (any, error) {
				return bind.URIParameters[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "uint: 10a is not a valid uint: invalid syntax",
		},
		{
			description: "WhenUint8IsInvalid_ShouldReturnError",
			parameters:  map[string]string{"uint8": "10a"},
			function: func(r *lit.Request) (any, error) {
				return bind.URIParameters[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "uint8: 10a is not a valid uint8: invalid syntax",
		},
		{
			description: "WhenUint16IsInvalid_ShouldReturnError",
			parameters:  map[string]string{"uint16": "10a"},
			function: func(r *lit.Request) (any, error) {
				return bind.URIParameters[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "uint16: 10a is not a valid uint16: invalid syntax",
		},
		{
			description: "WhenUint32IsInvalid_ShouldReturnError",
			parameters:  map[string]string{"uint32": "10a"},
			function: func(r *lit.Request) (any, error) {
				return bind.URIParameters[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "uint32: 10a is not a valid uint32: invalid syntax",
		},
		{
			description: "WhenUint64IsInvalid_ShouldReturnError",
			parameters:  map[string]string{"uint64": "10a"},
			function: func(r *lit.Request) (any, error) {
				return bind.URIParameters[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "uint64: 10a is not a valid uint64: invalid syntax",
		},
		{
			description: "WhenIntIsInvalid_ShouldReturnError",
			parameters:  map[string]string{"int": "10a"},
			function: func(r *lit.Request) (any, error) {
				return bind.URIParameters[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "int: 10a is not a valid int: invalid syntax",
		},
		{
			description: "WhenInt8IsInvalid_ShouldReturnError",
			parameters:  map[string]string{"int8": "10a"},
			function: func(r *lit.Request) (any, error) {
				return bind.URIParameters[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "int8: 10a is not a valid int8: invalid syntax",
		},
		{
			description: "WhenInt16IsInvalid_ShouldReturnError",
			parameters:  map[string]string{"int16": "10a"},
			function: func(r *lit.Request) (any, error) {
				return bind.URIParameters[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "int16: 10a is not a valid int16: invalid syntax",
		},
		{
			description: "WhenInt32IsInvalid_ShouldReturnError",
			parameters:  map[string]string{"int32": "10a"},
			function: func(r *lit.Request) (any, error) {
				return bind.URIParameters[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "int32: 10a is not a valid int32: invalid syntax",
		},
		{
			description: "WhenInt64IsInvalid_ShouldReturnError",
			parameters:  map[string]string{"int64": "10a"},
			function: func(r *lit.Request) (any, error) {
				return bind.URIParameters[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "int64: 10a is not a valid int64: invalid syntax",
		},
		{
			description: "WhenFloat32IsInvalid_ShouldReturnError",
			parameters:  map[string]string{"float32": "10a"},
			function: func(r *lit.Request) (any, error) {
				return bind.URIParameters[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "float32: 10a is not a valid float32: invalid syntax",
		},
		{
			description: "WhenFloat64IsInvalid_ShouldReturnError",
			parameters:  map[string]string{"float64": "10a"},
			function: func(r *lit.Request) (any, error) {
				return bind.URIParameters[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "float64: 10a is not a valid float64: invalid syntax",
		},
		{
			description: "WhenComplex64IsInvalid_ShouldReturnError",
			parameters:  map[string]string{"complex64": "10a"},
			function: func(r *lit.Request) (any, error) {
				return bind.URIParameters[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "complex64: 10a is not a valid complex64: invalid syntax",
		},
		{
			description: "WhenComplex128IsInvalid_ShouldReturnError",
			parameters:  map[string]string{"complex128": "10a"},
			function: func(r *lit.Request) (any, error) {
				return bind.URIParameters[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "complex128: 10a is not a valid complex128: invalid syntax",
		},
		{
			description: "WhenBoolIsInvalid_ShouldReturnError",
			parameters:  map[string]string{"bool": "10a"},
			function: func(r *lit.Request) (any, error) {
				return bind.URIParameters[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "bool: 10a is not a valid bool: invalid syntax",
		},
		{
			description: "WhenTimeIsInvalid_ShouldReturnError",
			parameters:  map[string]string{"time": "10a"},
			function: func(r *lit.Request) (any, error) {
				return bind.URIParameters[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError: `time: 10a is not a valid time.Time: parsing time "10a" as "2006-01-02T15:04:05Z07:00": ` +
				`cannot parse "10a" as "2006"`,
		},
		{
			description: "WhenTypeParameterIsValidatableWithValueReceiver_ShouldNotValidate",
			parameters:  map[string]string{"string": "string"},
			function: func(r *lit.Request) (any, error) {
				return bind.URIParameters[nonPointerReceiverValidatableFields](r)
			},
			expectedResult: nonPointerReceiverValidatableFields{String: "string"},
		},
		{
			description: "WhenTypeParameterIsValidatableWithPointerReceiver_ShouldValidate",
			parameters:  map[string]string{"string": "string"},
			function: func(r *lit.Request) (any, error) {
				return bind.URIParameters[pointerReceiverValidatableFields](r)
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
			request := httptest.NewRequest(http.MethodGet, "/", nil)

			r := lit.NewRequest(request).WithURIParameters(test.parameters)

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

func ExampleURIParameters() {
	r := lit.NewRequest(
		httptest.NewRequest(http.MethodGet, "/users/123/books/book_1", nil),
	).WithURIParameters(
		map[string]string{"user_id": "123", "book_id": "book_1"},
	)

	type RequestURIParameters struct {
		UserID int    `uri:"user_id"`
		BookID string `uri:"book_id"`
	}

	uri, err := bind.URIParameters[RequestURIParameters](r)
	if err == nil {
		fmt.Println(uri.UserID, uri.BookID)
	}

	// Output: 123 book_1
}
