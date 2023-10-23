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

func TestURLParameters_WhenTypeParameterIsNotStruct_ShouldPanic(t *testing.T) {
	t.Parallel()

	request := lit.NewRequest(
		httptest.NewRequest(http.MethodGet, "/", nil),
		nil,
	)

	require.PanicsWithValue(t, "T must be a struct type",
		func() { _, _ = bind.URLParameters[int](request) })
}

func TestURLParameters_WhenFieldHasUnsupportedType_ShouldPanic(t *testing.T) {
	t.Parallel()

	type fieldStruct struct {
		Field int
	}

	type targetStruct struct {
		Field fieldStruct `uri:"field"`
	}

	// Arrange
	parameters := map[string]string{
		"field": "123",
	}

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	request := lit.NewRequest(r, parameters)

	// Act
	// Assert
	require.PanicsWithValue(t, "unsupported type fieldStruct",
		func() { _, _ = bind.URLParameters[targetStruct](request) })
}

func TestURLParameters_ShouldBindSupportedTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description    string
		parameters     map[string]string
		expectedResult bindableTypes
		expectedError  string
	}{
		{
			description: "Valid",
			parameters: map[string]string{
				"string":     "hi",
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
				Slice:      nil,
				Array:      [2]int{},
			},
		},
		{
			description:   "Invalid uint",
			parameters:    map[string]string{"uint": "10a"},
			expectedError: "uint: 10a is not a valid uint: invalid syntax",
		},
		{
			description:   "Invalid uint8",
			parameters:    map[string]string{"uint8": "10a"},
			expectedError: "uint8: 10a is not a valid uint8: invalid syntax",
		},
		{
			description:   "Invalid uint16",
			parameters:    map[string]string{"uint16": "10a"},
			expectedError: "uint16: 10a is not a valid uint16: invalid syntax",
		},
		{
			description:   "Invalid uint32",
			parameters:    map[string]string{"uint32": "10a"},
			expectedError: "uint32: 10a is not a valid uint32: invalid syntax",
		},
		{
			description:   "Invalid uint64",
			parameters:    map[string]string{"uint64": "10a"},
			expectedError: "uint64: 10a is not a valid uint64: invalid syntax",
		},
		{
			description:   "Invalid int",
			parameters:    map[string]string{"int": "10a"},
			expectedError: "int: 10a is not a valid int: invalid syntax",
		},
		{
			description:   "Invalid int8",
			parameters:    map[string]string{"int8": "10a"},
			expectedError: "int8: 10a is not a valid int8: invalid syntax",
		},
		{
			description:   "Invalid int16",
			parameters:    map[string]string{"int16": "10a"},
			expectedError: "int16: 10a is not a valid int16: invalid syntax",
		},
		{
			description:   "Invalid int32",
			parameters:    map[string]string{"int32": "10a"},
			expectedError: "int32: 10a is not a valid int32: invalid syntax",
		},
		{
			description:   "Invalid int64",
			parameters:    map[string]string{"int64": "10a"},
			expectedError: "int64: 10a is not a valid int64: invalid syntax",
		},
		{
			description:   "Invalid float32",
			parameters:    map[string]string{"float32": "10a"},
			expectedError: "float32: 10a is not a valid float32: invalid syntax",
		},
		{
			description:   "Invalid float64",
			parameters:    map[string]string{"float64": "10a"},
			expectedError: "float64: 10a is not a valid float64: invalid syntax",
		},
		{
			description:   "Invalid complex64",
			parameters:    map[string]string{"complex64": "10a"},
			expectedError: "complex64: 10a is not a valid complex64: invalid syntax",
		},
		{
			description:   "Invalid complex128",
			parameters:    map[string]string{"complex128": "10a"},
			expectedError: "complex128: 10a is not a valid complex128: invalid syntax",
		},
		{
			description:   "Invalid bool",
			parameters:    map[string]string{"bool": "10a"},
			expectedError: "bool: 10a is not a valid bool: invalid syntax",
		},
		{
			description: "Invalid time",
			parameters:  map[string]string{"time": "10a"},
			expectedError: `time: 10a is not a valid time.Time: parsing time "10a" as "2006-01-02T15:04:05Z07:00": ` +
				`cannot parse "10a" as "2006"`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			request := lit.NewRequest(r, test.parameters)

			// Act
			result, err := bind.URLParameters[bindableTypes](request)

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

func TestURLParameters_WhenTagsAreNotPresentOrFieldIsUnexported_ShouldIgnore(t *testing.T) {
	t.Parallel()

	type targetStruct struct {
		ExportedAndPresent      string `uri:"exported"`
		ExportedAndNotPresent   string `uri:"not_present"`
		unexportedAndPresent    int    `uri:"unexported"`
		unexportedAndNotPresent int    `uri:"unexported_not_present"`
		Untagged                string
	}

	// Arrange
	parameters := map[string]string{
		"exported":   "123",
		"unexported": "123",
	}

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	request := lit.NewRequest(r, parameters)

	// Act
	result, err := bind.URLParameters[targetStruct](request)

	// Assert
	require.NoError(t, err)
	require.Equal(t, targetStruct{ExportedAndPresent: "123"}, result)
}
