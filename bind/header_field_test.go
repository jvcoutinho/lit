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

func TestHeaderField_ShouldBindSupportedTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description    string
		header         http.Header
		parameter      string
		function       func(*lit.Request, string) (any, error)
		expectedResult any
		expectedError  string
		shouldPanic    bool
	}{
		{
			description: "WhenHeaderHasNotBeenDefined_AndTypeParameterIsString_ShouldReturnEmptyString",
			header:      http.Header{"param": {"value"}},
			parameter:   "not_defined",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[string](r, p)
			},
			expectedResult: "",
		},
		{
			description: "ShouldBindString",
			header:      http.Header{"string": {"value"}},
			parameter:   "string",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[string](r, p)
			},
			expectedResult: "value",
		},
		{
			description: "WhenUintIsValid_ShouldBind",
			header:      http.Header{"uint": {"10"}},
			parameter:   "uint",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[uint](r, p)
			},
			expectedResult: uint(10),
		},
		{
			description: "WhenUintIsInvalid_ShouldReturnError",
			header:      http.Header{"uint": {"10a"}},
			parameter:   "uint",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[uint](r, p)
			},
			expectedResult: uint(0),
			expectedError:  "uint: 10a is not a valid uint: invalid syntax",
		},
		{
			description: "WhenUint8IsValid_ShouldBind",
			header:      http.Header{"uint8": {"10"}},
			parameter:   "uint8",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[uint8](r, p)
			},
			expectedResult: uint8(10),
		},
		{
			description: "WhenUint8IsInvalid_ShouldReturnError",
			header:      http.Header{"uint8": {"10a"}},
			parameter:   "uint8",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[uint8](r, p)
			},
			expectedResult: uint8(0),
			expectedError:  "uint8: 10a is not a valid uint8: invalid syntax",
		},
		{
			description: "WhenUint16IsValid_ShouldBind",
			header:      http.Header{"uint16": {"10"}},
			parameter:   "uint16",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[uint16](r, p)
			},
			expectedResult: uint16(10),
		},
		{
			description: "WhenUint16IsInvalid_ShouldReturnError",
			header:      http.Header{"uint16": {"10a"}},
			parameter:   "uint16",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[uint16](r, p)
			},
			expectedResult: uint16(0),
			expectedError:  "uint16: 10a is not a valid uint16: invalid syntax",
		},
		{
			description: "WhenUint32IsValid_ShouldBind",
			header:      http.Header{"uint32": {"10"}},
			parameter:   "uint32",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[uint32](r, p)
			},
			expectedResult: uint32(10),
		},
		{
			description: "WhenUint32IsInvalid_ShouldReturnError",
			header:      http.Header{"uint32": {"10a"}},
			parameter:   "uint32",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[uint32](r, p)
			},
			expectedResult: uint32(0),
			expectedError:  "uint32: 10a is not a valid uint32: invalid syntax",
		},
		{
			description: "WhenUint64IsValid_ShouldBind",
			header:      http.Header{"uint64": {"10"}},
			parameter:   "uint64",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[uint64](r, p)
			},
			expectedResult: uint64(10),
		},
		{
			description: "WhenUint64IsInvalid_ShouldReturnError",
			header:      http.Header{"uint64": {"10a"}},
			parameter:   "uint64",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[uint64](r, p)
			},
			expectedResult: uint64(0),
			expectedError:  "uint64: 10a is not a valid uint64: invalid syntax",
		},
		{
			description: "WhenIntIsValid_ShouldBind",
			header:      http.Header{"int": {"10"}},
			parameter:   "int",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[int](r, p)
			},
			expectedResult: 10,
		},
		{
			description: "WhenIntIsInvalid_ShouldReturnError",
			header:      http.Header{"int": {"10a"}},
			parameter:   "int",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[int](r, p)
			},
			expectedResult: 0,
			expectedError:  "int: 10a is not a valid int: invalid syntax",
		},
		{
			description: "WhenInt8IsValid_ShouldBind",
			header:      http.Header{"int8": {"10"}},
			parameter:   "int8",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[int8](r, p)
			},
			expectedResult: int8(10),
		},
		{
			description: "WhenInt8IsInvalid_ShouldReturnError",
			header:      http.Header{"int8": {"10a"}},
			parameter:   "int8",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[int8](r, p)
			},
			expectedResult: int8(0),
			expectedError:  "int8: 10a is not a valid int8: invalid syntax",
		},
		{
			description: "WhenInt16IsValid_ShouldBind",
			header:      http.Header{"int16": {"10"}},
			parameter:   "int16",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[int16](r, p)
			},
			expectedResult: int16(10),
		},
		{
			description: "WhenInt16IsInvalid_ShouldReturnError",
			header:      http.Header{"int16": {"10a"}},
			parameter:   "int16",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[int16](r, p)
			},
			expectedResult: int16(0),
			expectedError:  "int16: 10a is not a valid int16: invalid syntax",
		},
		{
			description: "WhenInt32IsValid_ShouldBind",
			header:      http.Header{"int32": {"10"}},
			parameter:   "int32",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[int32](r, p)
			},
			expectedResult: int32(10),
		},
		{
			description: "WhenInt32IsInvalid_ShouldReturnError",
			header:      http.Header{"int32": {"10a"}},
			parameter:   "int32",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[int32](r, p)
			},
			expectedResult: int32(0),
			expectedError:  "int32: 10a is not a valid int32: invalid syntax",
		},
		{
			description: "WhenInt64IsValid_ShouldBind",
			header:      http.Header{"int64": {"10"}},
			parameter:   "int64",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[int64](r, p)
			},
			expectedResult: int64(10),
		},
		{
			description: "WhenInt64IsInvalid_ShouldReturnError",
			header:      http.Header{"int64": {"10a"}},
			parameter:   "int64",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[int64](r, p)
			},
			expectedResult: int64(0),
			expectedError:  "int64: 10a is not a valid int64: invalid syntax",
		},
		{
			description: "WhenFloat32IsValid_ShouldBind",
			header:      http.Header{"float32": {"10"}},
			parameter:   "float32",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[float32](r, p)
			},
			expectedResult: float32(10.0),
		},
		{
			description: "WhenFloat32IsInvalid_ShouldReturnError",
			header:      http.Header{"float32": {"10a"}},
			parameter:   "float32",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[float32](r, p)
			},
			expectedResult: float32(0),
			expectedError:  "float32: 10a is not a valid float32: invalid syntax",
		},
		{
			description: "WhenFloat64IsValid_ShouldBind",
			header:      http.Header{"float64": {"10"}},
			parameter:   "float64",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[float64](r, p)
			},
			expectedResult: 10.0,
		},
		{
			description: "WhenFloat64IsInvalid_ShouldReturnError",
			header:      http.Header{"float64": {"10a"}},
			parameter:   "float64",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[float64](r, p)
			},
			expectedResult: 0.0,
			expectedError:  "float64: 10a is not a valid float64: invalid syntax",
		},
		{
			description: "WhenComplex64IsValid_ShouldBind",
			header:      http.Header{"complex64": {"10"}},
			parameter:   "complex64",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[complex64](r, p)
			},
			expectedResult: complex64(10 + 0i),
		},
		{
			description: "WhenComplex64IsInvalid_ShouldReturnError",
			header:      http.Header{"complex64": {"10a"}},
			parameter:   "complex64",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[complex64](r, p)
			},
			expectedResult: complex64(0),
			expectedError:  "complex64: 10a is not a valid complex64: invalid syntax",
		},
		{
			description: "WhenComplex128IsValid_ShouldBind",
			header:      http.Header{"complex128": {"10"}},
			parameter:   "complex128",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[complex128](r, p)
			},
			expectedResult: 10 + 0i,
		},
		{
			description: "WhenComplex128IsInvalid_ShouldReturnError",
			header:      http.Header{"complex128": {"10a"}},
			parameter:   "complex128",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[complex128](r, p)
			},
			expectedResult: 0 + 0i,
			expectedError:  "complex128: 10a is not a valid complex128: invalid syntax",
		},
		{
			description: "WhenBoolIsValid_ShouldBind",
			header:      http.Header{"bool": {"true"}},
			parameter:   "bool",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[bool](r, p)
			},
			expectedResult: true,
		},
		{
			description: "WhenBoolIsInvalid_ShouldReturnError",
			header:      http.Header{"bool": {"10a"}},
			parameter:   "bool",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[bool](r, p)
			},
			expectedResult: false,
			expectedError:  "bool: 10a is not a valid bool: invalid syntax",
		},
		{
			description: "WhenTimeIsValid_ShouldBind",
			header:      http.Header{"time": {"2023-10-22T00:00:00Z"}},
			parameter:   "time",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[time.Time](r, p)
			},
			expectedResult: time.Date(2023, 10, 22, 0, 0, 0, 0, time.UTC),
		},
		{
			description: "WhenTimeIsInvalid_ShouldReturnError",
			header:      http.Header{"time": {"10a"}},
			parameter:   "time",
			function: func(r *lit.Request, p string) (any, error) {
				return bind.HeaderField[time.Time](r, p)
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
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			for key, value := range test.header {
				for _, v := range value {
					request.Header.Add(key, v)
				}
			}

			r := lit.NewRequest(request, nil)

			// Act
			if test.shouldPanic {
				require.PanicsWithValue(t, test.expectedError, func() {
					_, _ = test.function(r, test.parameter)
				})

				return
			}

			result, err := test.function(r, test.parameter)

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

func ExampleHeaderField() {
	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	req.Header.Add("Content-Length", "150")
	req.Header.Add("Authorization", "Bearer uPSsoa65gqkFv2Z6sZ3rZCZwnCjzaXe8TNdk0bJCFFJGrH6wmnzyK4evHBtTuvVH")

	r := lit.NewRequest(req, nil)

	contentLength, err := bind.HeaderField[int](r, "Content-Length")
	if err == nil {
		fmt.Println(contentLength)
	}

	authorization, err := bind.HeaderField[string](r, "authorization") // case-insensitive
	if err == nil {
		fmt.Println(authorization)
	}

	// Output:
	// 150
	// Bearer uPSsoa65gqkFv2Z6sZ3rZCZwnCjzaXe8TNdk0bJCFFJGrH6wmnzyK4evHBtTuvVH
}
