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

func TestURIParameter_WhenParameterIsNotDefined_ShouldPanic(t *testing.T) {
	t.Parallel()

	// Arrange
	request := lit.NewRequest(
		httptest.NewRequest(http.MethodGet, "/", nil),
		map[string]string{"user_id": "123"},
	)

	// Act
	// Assert
	require.PanicsWithValue(t, "book_id has not been defined as one of the request parameters: [user_id]", func() {
		_, _ = bind.URIParameter[int](request, "book_id")
	})
}

func TestURIParameter_ShouldBindSupportedTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description    string
		parameters     map[string]string
		parameter      string
		bindFunction   func(*lit.Request, string) (any, error)
		expectedResult any
		expectedError  string
	}{
		{
			description: "Valid string",
			parameters:  map[string]string{"string": "hi"},
			parameter:   "string",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[string](r, p)
			},
			expectedResult: "hi",
			expectedError:  "",
		},
		{
			description: "Valid uint",
			parameters:  map[string]string{"uint": "10"},
			parameter:   "uint",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[uint](r, p)
			},
			expectedResult: uint(10),
			expectedError:  "",
		},
		{
			description: "Invalid uint",
			parameters:  map[string]string{"uint": "10a"},
			parameter:   "uint",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[uint](r, p)
			},
			expectedResult: uint(0),
			expectedError:  "uint: 10a is not a valid uint: invalid syntax",
		},
		{
			description: "Valid uint8",
			parameters:  map[string]string{"uint8": "10"},
			parameter:   "uint8",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[uint8](r, p)
			},
			expectedResult: uint8(10),
			expectedError:  "",
		},
		{
			description: "Invalid uint8",
			parameters:  map[string]string{"uint8": "10a"},
			parameter:   "uint8",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[uint8](r, p)
			},
			expectedResult: uint8(0),
			expectedError:  "uint8: 10a is not a valid uint8: invalid syntax",
		},
		{
			description: "Valid uint16",
			parameters:  map[string]string{"uint16": "10"},
			parameter:   "uint16",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[uint16](r, p)
			},
			expectedResult: uint16(10),
			expectedError:  "",
		},
		{
			description: "Invalid uint16",
			parameters:  map[string]string{"uint16": "10a"},
			parameter:   "uint16",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[uint16](r, p)
			},
			expectedResult: uint16(0),
			expectedError:  "uint16: 10a is not a valid uint16: invalid syntax",
		},
		{
			description: "Valid uint32",
			parameters:  map[string]string{"uint32": "10"},
			parameter:   "uint32",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[uint32](r, p)
			},
			expectedResult: uint32(10),
			expectedError:  "",
		},
		{
			description: "Invalid uint32",
			parameters:  map[string]string{"uint32": "10a"},
			parameter:   "uint32",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[uint32](r, p)
			},
			expectedResult: uint32(0),
			expectedError:  "uint32: 10a is not a valid uint32: invalid syntax",
		},
		{
			description: "Valid uint64",
			parameters:  map[string]string{"uint64": "10"},
			parameter:   "uint64",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[uint64](r, p)
			},
			expectedResult: uint64(10),
			expectedError:  "",
		},
		{
			description: "Invalid uint64",
			parameters:  map[string]string{"uint64": "10a"},
			parameter:   "uint64",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[uint64](r, p)
			},
			expectedResult: uint64(0),
			expectedError:  "uint64: 10a is not a valid uint64: invalid syntax",
		},
		{
			description: "Valid int",
			parameters:  map[string]string{"int": "10"},
			parameter:   "int",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[int](r, p)
			},
			expectedResult: 10,
			expectedError:  "",
		},
		{
			description: "Invalid int",
			parameters:  map[string]string{"int": "10a"},
			parameter:   "int",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[int](r, p)
			},
			expectedResult: 0,
			expectedError:  "int: 10a is not a valid int: invalid syntax",
		},
		{
			description: "Valid int8",
			parameters:  map[string]string{"int8": "10"},
			parameter:   "int8",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[int8](r, p)
			},
			expectedResult: int8(10),
			expectedError:  "",
		},
		{
			description: "Invalid int8",
			parameters:  map[string]string{"int8": "10a"},
			parameter:   "int8",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[int8](r, p)
			},
			expectedResult: int8(0),
			expectedError:  "int8: 10a is not a valid int8: invalid syntax",
		},
		{
			description: "Valid int16",
			parameters:  map[string]string{"int16": "10"},
			parameter:   "int16",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[int16](r, p)
			},
			expectedResult: int16(10),
			expectedError:  "",
		},
		{
			description: "Invalid int16",
			parameters:  map[string]string{"int16": "10a"},
			parameter:   "int16",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[int16](r, p)
			},
			expectedResult: int16(0),
			expectedError:  "int16: 10a is not a valid int16: invalid syntax",
		},
		{
			description: "Valid int32",
			parameters:  map[string]string{"int32": "10"},
			parameter:   "int32",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[int32](r, p)
			},
			expectedResult: int32(10),
			expectedError:  "",
		},
		{
			description: "Invalid int32",
			parameters:  map[string]string{"int32": "10a"},
			parameter:   "int32",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[int32](r, p)
			},
			expectedResult: int32(0),
			expectedError:  "int32: 10a is not a valid int32: invalid syntax",
		},
		{
			description: "Valid int64",
			parameters:  map[string]string{"int64": "10"},
			parameter:   "int64",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[int64](r, p)
			},
			expectedResult: int64(10),
			expectedError:  "",
		},
		{
			description: "Invalid int64",
			parameters:  map[string]string{"int64": "10a"},
			parameter:   "int64",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[int64](r, p)
			},
			expectedResult: int64(0),
			expectedError:  "int64: 10a is not a valid int64: invalid syntax",
		},
		{
			description: "Valid float32",
			parameters:  map[string]string{"float32": "10"},
			parameter:   "float32",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[float32](r, p)
			},
			expectedResult: float32(10.0),
			expectedError:  "",
		},
		{
			description: "Invalid float32",
			parameters:  map[string]string{"float32": "10a"},
			parameter:   "float32",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[float32](r, p)
			},
			expectedResult: float32(0),
			expectedError:  "float32: 10a is not a valid float32: invalid syntax",
		},
		{
			description: "Valid float64",
			parameters:  map[string]string{"float64": "10"},
			parameter:   "float64",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[float64](r, p)
			},
			expectedResult: 10.0,
			expectedError:  "",
		},
		{
			description: "Invalid float64",
			parameters:  map[string]string{"float64": "10a"},
			parameter:   "float64",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[float64](r, p)
			},
			expectedResult: 0.0,
			expectedError:  "float64: 10a is not a valid float64: invalid syntax",
		},
		{
			description: "Valid complex64",
			parameters:  map[string]string{"complex64": "10"},
			parameter:   "complex64",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[complex64](r, p)
			},
			expectedResult: complex64(10 + 0i),
			expectedError:  "",
		},
		{
			description: "Invalid complex64",
			parameters:  map[string]string{"complex64": "10a"},
			parameter:   "complex64",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[complex64](r, p)
			},
			expectedResult: complex64(0),
			expectedError:  "complex64: 10a is not a valid complex64: invalid syntax",
		},
		{
			description: "Valid bool",
			parameters:  map[string]string{"bool": "true"},
			parameter:   "bool",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[bool](r, p)
			},
			expectedResult: true,
			expectedError:  "",
		},
		{
			description: "Invalid bool",
			parameters:  map[string]string{"bool": "10a"},
			parameter:   "bool",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[bool](r, p)
			},
			expectedResult: false,
			expectedError:  "bool: 10a is not a valid bool: invalid syntax",
		},
		{
			description: "Valid time",
			parameters:  map[string]string{"time": "2023-10-22T00:00:00Z"},
			parameter:   "time",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[time.Time](r, p)
			},
			expectedResult: time.Date(2023, 10, 22, 0, 0, 0, 0, time.UTC),
			expectedError:  "",
		},
		{
			description: "Invalid time",
			parameters:  map[string]string{"time": "10a"},
			parameter:   "time",
			bindFunction: func(r *lit.Request, p string) (any, error) {
				return bind.URIParameter[time.Time](r, p)
			},
			expectedResult: time.Time{},
			expectedError: `time: 10a is not a valid time.Time: parsing time "10a" as "2006-01-02T15:04:05Z07:00": ` +
				`cannot parse "10a" as "2006"`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			request := lit.NewRequest(
				httptest.NewRequest(http.MethodGet, "/", nil),
				test.parameters,
			)

			// Act
			result, err := test.bindFunction(request, test.parameter)

			// Assert
			if test.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, test.expectedError)
			}

			require.Equal(t, test.expectedResult, result)
		})
	}
}

func ExampleURIParameter() {
	r := lit.NewRequest(
		httptest.NewRequest(http.MethodGet, "/users/123/books/book_1", nil),
		map[string]string{"user_id": "123", "book_id": "book_1"},
	)

	userID, err := bind.URIParameter[int](r, "user_id")
	if err == nil {
		fmt.Println(userID)
	}

	// Output: 123
}
