package bind_test

import (
	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/bind"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHeader_WhenTypeParameterIsNotStruct_ShouldPanic(t *testing.T) {
	t.Parallel()

	request := lit.NewRequest(
		httptest.NewRequest(http.MethodGet, "/", nil),
		nil,
	)

	require.PanicsWithValue(t, "T must be a struct type",
		func() { _, _ = bind.Header[int](request) })
}

func TestHeader_WhenFieldHasUnsupportedType_ShouldPanic(t *testing.T) {
	t.Parallel()

	type fieldStruct struct {
		Field int
	}

	type targetStruct struct {
		Field fieldStruct `header:"field"`
	}

	// Arrange
	header := http.Header{
		"field": {"123"},
	}

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header = header

	request := lit.NewRequest(r, nil)

	// Act
	// Assert
	require.PanicsWithValue(t, "unsupported type fieldStruct",
		func() { _, _ = bind.Header[targetStruct](request) })
}

func TestHeader_ShouldBindSupportedTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description    string
		header         http.Header
		expectedResult bindableTypes
		expectedError  string
	}{
		{
			description: "Valid",
			header: map[string][]string{
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
			expectedResult: bindableTypes{
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
			description:   "Invalid uint",
			header:        map[string][]string{"uint": {"10a"}},
			expectedError: "uint: 10a is not a valid uint: invalid syntax",
		},
		{
			description:   "Invalid uint8",
			header:        map[string][]string{"uint8": {"10a"}},
			expectedError: "uint8: 10a is not a valid uint8: invalid syntax",
		},
		{
			description:   "Invalid uint16",
			header:        map[string][]string{"uint16": {"10a"}},
			expectedError: "uint16: 10a is not a valid uint16: invalid syntax",
		},
		{
			description:   "Invalid uint32",
			header:        map[string][]string{"uint32": {"10a"}},
			expectedError: "uint32: 10a is not a valid uint32: invalid syntax",
		},
		{
			description:   "Invalid uint64",
			header:        map[string][]string{"uint64": {"10a"}},
			expectedError: "uint64: 10a is not a valid uint64: invalid syntax",
		},
		{
			description:   "Invalid int",
			header:        map[string][]string{"int": {"10a"}},
			expectedError: "int: 10a is not a valid int: invalid syntax",
		},
		{
			description:   "Invalid int8",
			header:        map[string][]string{"int8": {"10a"}},
			expectedError: "int8: 10a is not a valid int8: invalid syntax",
		},
		{
			description:   "Invalid int16",
			header:        map[string][]string{"int16": {"10a"}},
			expectedError: "int16: 10a is not a valid int16: invalid syntax",
		},
		{
			description:   "Invalid int32",
			header:        map[string][]string{"int32": {"10a"}},
			expectedError: "int32: 10a is not a valid int32: invalid syntax",
		},
		{
			description:   "Invalid int64",
			header:        map[string][]string{"int64": {"10a"}},
			expectedError: "int64: 10a is not a valid int64: invalid syntax",
		},
		{
			description:   "Invalid float32",
			header:        map[string][]string{"float32": {"10a"}},
			expectedError: "float32: 10a is not a valid float32: invalid syntax",
		},
		{
			description:   "Invalid float64",
			header:        map[string][]string{"float64": {"10a"}},
			expectedError: "float64: 10a is not a valid float64: invalid syntax",
		},
		{
			description:   "Invalid complex64",
			header:        map[string][]string{"complex64": {"10a"}},
			expectedError: "complex64: 10a is not a valid complex64: invalid syntax",
		},
		{
			description:   "Invalid complex128",
			header:        map[string][]string{"complex128": {"10a"}},
			expectedError: "complex128: 10a is not a valid complex128: invalid syntax",
		},
		{
			description:   "Invalid bool",
			header:        map[string][]string{"bool": {"10a"}},
			expectedError: "bool: 10a is not a valid bool: invalid syntax",
		},
		{
			description: "Invalid time",
			header:      map[string][]string{"time": {"10a"}},
			expectedError: `time: 10a is not a valid time.Time: parsing time "10a" as "2006-01-02T15:04:05Z07:00": ` +
				`cannot parse "10a" as "2006"`,
		},
		{
			description:   "Invalid slice element",
			header:        map[string][]string{"slice": {"10a"}},
			expectedError: "slice: 10a is not a valid int: invalid syntax",
		},
		{
			description:   "Invalid array element",
			header:        map[string][]string{"array": {"10a"}},
			expectedError: "array: 10a is not a valid int: invalid syntax",
		},
		{
			description:   "Invalid array length",
			header:        map[string][]string{"array": {"10", "20", "30"}},
			expectedError: "array: expected at most 2 elements. Got 3",
		},
		{
			description:   "Invalid field not slice",
			header:        map[string][]string{"int": {"10", "20"}},
			expectedError: "int: [10 20] is not a valid int",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.Header = test.header

			request := lit.NewRequest(r, nil)

			// Act
			result, err := bind.Header[bindableTypes](request)

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

func TestHeader_WhenTagsAreNotPresentOrFieldIsUnexported_ShouldIgnore(t *testing.T) {
	t.Parallel()

	type targetStruct struct {
		ExportedAndPresent      string `header:"exported"`
		ExportedAndNotPresent   string `header:"not_present"`
		unexportedAndPresent    int    `header:"unexported"`
		unexportedAndNotPresent int    `header:"unexported_not_present"`
		Untagged                string
	}

	// Arrange
	header := http.Header{
		"exported":   {"123"},
		"unexported": {"123"},
	}

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header = header

	request := lit.NewRequest(r, nil)

	// Act
	result, err := bind.Header[targetStruct](request)

	// Assert
	require.NoError(t, err)
	require.Equal(t, targetStruct{ExportedAndPresent: "123"}, result)
}
