package bind_test

import (
	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/bind"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
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

func TestHeader(t *testing.T) {
	t.Parallel()

	type targetStruct struct {
		ContentType   string `header:"Content-Type"`
		Authorization string `header:"Authorization"`
		ContentLength int    `header:"Content-Length"`
		Accept        string
	}

	tests := []struct {
		description    string
		header         http.Header
		expectedResult any
		expectedError  string
	}{
		{
			description: "WhenArgumentsMatchBindingTarget_ShouldBind",
			header: map[string][]string{
				"Content-Type":   {"application/json"},
				"Authorization":  {"Bearer token"},
				"Content-Length": {"32"},
			},
			expectedResult: targetStruct{
				ContentType:   "application/json",
				Authorization: "Bearer token",
				ContentLength: 32,
			},
			expectedError: "",
		},
		{
			description: "WhenHeadersAreMissingInRequest_ShouldIgnoreThem",
			header: map[string][]string{
				"Content-Type": {"application/json"},
			},
			expectedResult: targetStruct{
				ContentType: "application/json",
			},
			expectedError: "",
		},
		{
			description: "WhenHeadersDoNotMatchBindingTarget_ShouldReturnError",
			header: map[string][]string{
				"Content-Length": {"123a"},
			},
			expectedResult: targetStruct{},
			expectedError:  "Content-Length: 123a is not a valid int: invalid syntax",
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
			result, err := bind.Header[targetStruct](request)

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
