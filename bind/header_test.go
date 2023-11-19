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

func TestHeader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description    string
		header         http.Header
		function       func(r *lit.Request) (any, error)
		expectedResult any
		expectedError  string
		shouldPanic    bool
	}{
		{
			description: "WhenTypeParameterIsNotStruct_ShouldPanic",
			function: func(r *lit.Request) (any, error) {
				return bind.Header[int](r)
			},
			expectedError: "T must be a struct type",
			shouldPanic:   true,
		},
		{
			description: "WhenTargetHasUnbindableField_ShouldPanic",
			header: http.Header{
				"field": {"123"},
			},
			function: func(r *lit.Request) (any, error) {
				return bind.Header[unbindableField](r)
			},
			expectedError: "unbindable type lit.Request",
			shouldPanic:   true,
		},
		{
			description: "WhenFieldIsUnexported_ShouldIgnore",
			header: http.Header{
				"unexported": {"123"},
			},
			function: func(r *lit.Request) (any, error) {
				return bind.Header[ignorableFields](r)
			},
			expectedResult: ignorableFields{},
		},
		{
			description: "WhenFieldIsMissingFromRequest_ShouldIgnore",
			header:      http.Header{},
			function: func(r *lit.Request) (any, error) {
				return bind.Header[ignorableFields](r)
			},
			expectedResult: ignorableFields{},
		},
		{
			description: "WhenFieldsAreValid_ShouldBindThem",
			header: http.Header{
				"string":     {"hi"},
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
				return bind.Header[bindableFields](r)
			},
			expectedResult: bindableFields{
				String:     "hi",
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
			description: "WhenUintIsInvalid_ShouldReturnError",
			header:      http.Header{"uint": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Header[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "uint: 10a is not a valid uint: invalid syntax",
		},
		{
			description: "WhenUint8IsInvalid_ShouldReturnError",
			header:      http.Header{"uint8": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Header[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "uint8: 10a is not a valid uint8: invalid syntax",
		},
		{
			description: "WhenUint16IsInvalid_ShouldReturnError",
			header:      http.Header{"uint16": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Header[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "uint16: 10a is not a valid uint16: invalid syntax",
		},
		{
			description: "WhenUint32IsInvalid_ShouldReturnError",
			header:      http.Header{"uint32": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Header[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "uint32: 10a is not a valid uint32: invalid syntax",
		},
		{
			description: "WhenUint64IsInvalid_ShouldReturnError",
			header:      http.Header{"uint64": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Header[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "uint64: 10a is not a valid uint64: invalid syntax",
		},
		{
			description: "WhenIntIsInvalid_ShouldReturnError",
			header:      http.Header{"int": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Header[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "int: 10a is not a valid int: invalid syntax",
		},
		{
			description: "WhenInt8IsInvalid_ShouldReturnError",
			header:      http.Header{"int8": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Header[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "int8: 10a is not a valid int8: invalid syntax",
		},
		{
			description: "WhenInt16IsInvalid_ShouldReturnError",
			header:      http.Header{"int16": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Header[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "int16: 10a is not a valid int16: invalid syntax",
		},
		{
			description: "WhenInt32IsInvalid_ShouldReturnError",
			header:      http.Header{"int32": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Header[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "int32: 10a is not a valid int32: invalid syntax",
		},
		{
			description: "WhenInt64IsInvalid_ShouldReturnError",
			header:      http.Header{"int64": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Header[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "int64: 10a is not a valid int64: invalid syntax",
		},
		{
			description: "WhenFloat32IsInvalid_ShouldReturnError",
			header:      http.Header{"float32": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Header[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "float32: 10a is not a valid float32: invalid syntax",
		},
		{
			description: "WhenFloat64IsInvalid_ShouldReturnError",
			header:      http.Header{"float64": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Header[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "float64: 10a is not a valid float64: invalid syntax",
		},
		{
			description: "WhenComplex64IsInvalid_ShouldReturnError",
			header:      http.Header{"complex64": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Header[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "complex64: 10a is not a valid complex64: invalid syntax",
		},
		{
			description: "WhenComplex128IsInvalid_ShouldReturnError",
			header:      http.Header{"complex128": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Header[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "complex128: 10a is not a valid complex128: invalid syntax",
		},
		{
			description: "WhenBoolIsInvalid_ShouldReturnError",
			header:      http.Header{"bool": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Header[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "bool: 10a is not a valid bool: invalid syntax",
		},
		{
			description: "WhenTimeIsInvalid_ShouldReturnError",
			header:      http.Header{"time": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Header[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError: `time: 10a is not a valid time.Time: parsing time "10a" as "2006-01-02T15:04:05Z07:00": ` +
				`cannot parse "10a" as "2006"`,
		},
		{
			description: "WhenSliceElementIsInvalid_ShouldReturnError",
			header:      http.Header{"slice": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Header[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "slice: 10a is not a valid int: invalid syntax",
		},
		{
			description: "WhenArrayElementIsInvalid_ShouldReturnError",
			header:      http.Header{"array": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Header[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "array: 10a is not a valid int: invalid syntax",
		},
		{
			description: "WhenArrayLengthIsGreaterThanCapacity_ShouldReturnError",
			header:      http.Header{"array": {"10", "20", "30"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Header[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "array: [10 20 30] is not a valid [2]int: expected at most 2 elements. Got 3",
		},
		{
			description: "WhenFieldIsInvalidSlice_ShouldReturnError",
			header:      http.Header{"int": {"10", "20"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Header[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "int: [10 20] is not a valid int",
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

			require.Equal(t, test.expectedResult, result)
			require.Equal(t, test.expectedError, errMessage)
		})
	}
}

func ExampleHeader() {
	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	req.Header.Add("Content-Length", "150")
	req.Header.Add("Authorization", "Bearer uPSsoa65gqkFv2Z6sZ3rZCZwnCjzaXe8TNdk0bJCFFJGrH6wmnzyK4evHBtTuvVH")

	r := lit.NewRequest(req, nil)

	type Header struct {
		ContentLength uint   `header:"Content-Length"`
		Authorization string `header:"Authorization"`
	}

	h, err := bind.Header[Header](r)
	if err == nil {
		fmt.Println(h.ContentLength)
		fmt.Println(h.Authorization)
	}

	// Output:
	// 150
	// Bearer uPSsoa65gqkFv2Z6sZ3rZCZwnCjzaXe8TNdk0bJCFFJGrH6wmnzyK4evHBtTuvVH
}
