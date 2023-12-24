package bind_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/bind"
	"github.com/stretchr/testify/require"
)

func TestQuery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description    string
		query          url.Values
		function       func(r *lit.Request) (any, error)
		expectedResult any
		expectedError  string
		shouldPanic    bool
	}{
		{
			description: "WhenTypeParameterIsNotStruct_ShouldPanic",
			function: func(r *lit.Request) (any, error) {
				return bind.Query[int](r)
			},
			expectedError: "T must be a struct type",
			shouldPanic:   true,
		},
		{
			description: "WhenTargetHasUnbindableField_ShouldPanic",
			query: url.Values{
				"field": {"123"},
			},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[unbindableField](r)
			},
			expectedError: "unbindable type lit.Request",
			shouldPanic:   true,
		},
		{
			description: "WhenFieldIsUnexported_ShouldIgnore",
			query: url.Values{
				"unexported": {"123"},
			},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[ignorableFields](r)
			},
			expectedResult: ignorableFields{},
		},
		{
			description: "WhenFieldIsMissingFromRequest_ShouldIgnore",
			query:       url.Values{},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[ignorableFields](r)
			},
			expectedResult: ignorableFields{},
		},
		{
			description: "WhenFieldsAreValid_ShouldBindThem",
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
				return bind.Query[bindableFields](r)
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
			description: "WhenPointerValueIsInvalid_ShouldReturnError",
			query:       url.Values{"pointer": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "pointer: 10a is not a valid int: invalid syntax",
		},
		{
			description: "WhenUintIsInvalid_ShouldReturnError",
			query:       url.Values{"uint": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "uint: 10a is not a valid uint: invalid syntax",
		},
		{
			description: "WhenUint8IsInvalid_ShouldReturnError",
			query:       url.Values{"uint8": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "uint8: 10a is not a valid uint8: invalid syntax",
		},
		{
			description: "WhenUint16IsInvalid_ShouldReturnError",
			query:       url.Values{"uint16": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "uint16: 10a is not a valid uint16: invalid syntax",
		},
		{
			description: "WhenUint32IsInvalid_ShouldReturnError",
			query:       url.Values{"uint32": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "uint32: 10a is not a valid uint32: invalid syntax",
		},
		{
			description: "WhenUint64IsInvalid_ShouldReturnError",
			query:       url.Values{"uint64": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "uint64: 10a is not a valid uint64: invalid syntax",
		},
		{
			description: "WhenIntIsInvalid_ShouldReturnError",
			query:       url.Values{"int": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "int: 10a is not a valid int: invalid syntax",
		},
		{
			description: "WhenInt8IsInvalid_ShouldReturnError",
			query:       url.Values{"int8": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "int8: 10a is not a valid int8: invalid syntax",
		},
		{
			description: "WhenInt16IsInvalid_ShouldReturnError",
			query:       url.Values{"int16": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "int16: 10a is not a valid int16: invalid syntax",
		},
		{
			description: "WhenInt32IsInvalid_ShouldReturnError",
			query:       url.Values{"int32": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "int32: 10a is not a valid int32: invalid syntax",
		},
		{
			description: "WhenInt64IsInvalid_ShouldReturnError",
			query:       url.Values{"int64": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "int64: 10a is not a valid int64: invalid syntax",
		},
		{
			description: "WhenFloat32IsInvalid_ShouldReturnError",
			query:       url.Values{"float32": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "float32: 10a is not a valid float32: invalid syntax",
		},
		{
			description: "WhenFloat64IsInvalid_ShouldReturnError",
			query:       url.Values{"float64": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "float64: 10a is not a valid float64: invalid syntax",
		},
		{
			description: "WhenComplex64IsInvalid_ShouldReturnError",
			query:       url.Values{"complex64": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "complex64: 10a is not a valid complex64: invalid syntax",
		},
		{
			description: "WhenComplex128IsInvalid_ShouldReturnError",
			query:       url.Values{"complex128": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "complex128: 10a is not a valid complex128: invalid syntax",
		},
		{
			description: "WhenBoolIsInvalid_ShouldReturnError",
			query:       url.Values{"bool": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "bool: 10a is not a valid bool: invalid syntax",
		},
		{
			description: "WhenTimeIsInvalid_ShouldReturnError",
			query:       url.Values{"time": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError: `time: 10a is not a valid time.Time: parsing time "10a" as "2006-01-02T15:04:05Z07:00": ` +
				`cannot parse "10a" as "2006"`,
		},
		{
			description: "WhenSliceElementIsInvalid_ShouldReturnError",
			query:       url.Values{"slice": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "slice: 10a is not a valid int: invalid syntax",
		},
		{
			description: "WhenArrayElementIsInvalid_ShouldReturnError",
			query:       url.Values{"array": {"10a"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "array: 10a is not a valid int: invalid syntax",
		},
		{
			description: "WhenArrayLengthIsGreaterThanCapacity_ShouldReturnError",
			query:       url.Values{"array": {"10", "20", "30"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "array: [10 20 30] is not a valid [2]int: expected at most 2 elements. Got 3",
		},
		{
			description: "WhenFieldIsInvalidSlice_ShouldReturnError",
			query:       url.Values{"int": {"10", "20"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "int: [10 20] is not a valid int",
		},
		{
			description: "WhenTypeParameterIsValidatableWithValueReceiver_ShouldNotValidate",
			query:       map[string][]string{"string": {"string"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[nonPointerReceiverValidatableFields](r)
			},
			expectedResult: nonPointerReceiverValidatableFields{String: "string"},
		},
		{
			description: "WhenTypeParameterIsValidatableWithPointerReceiver_ShouldValidate",
			query:       map[string][]string{"string": {"string"}},
			function: func(r *lit.Request) (any, error) {
				return bind.Query[pointerReceiverValidatableFields](r)
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
			request.URL.RawQuery = test.query.Encode()

			r := lit.NewRequest(request)

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

func ExampleQuery() {
	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	req.URL.RawQuery = "publish_year=2009&name=Percy%20Jackson"

	r := lit.NewRequest(req)

	type BookQuery struct {
		PublishYear uint   `query:"publish_year"`
		Name        string `query:"name"`
	}

	query, err := bind.Query[BookQuery](r)
	if err == nil {
		fmt.Println(query.PublishYear, query.Name)
	}

	// Output: 2009 Percy Jackson
}
