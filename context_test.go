package lit_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jvcoutinho/lit"
	"github.com/stretchr/testify/require"
)

func TestContext_SetStatusCode(t *testing.T) {
	t.Parallel()

	// Arrange
	recorder := httptest.NewRecorder()
	context := lit.NewContext(recorder, nil)

	require.NotEqual(t, http.StatusBadGateway, recorder.Code)

	// Act
	context.SetStatusCode(http.StatusBadGateway)

	// Assert
	require.Equal(t, http.StatusBadGateway, recorder.Code)
}

func TestContext_WriteBody(t *testing.T) {
	t.Parallel()

	// Arrange
	recorder := httptest.NewRecorder()
	context := lit.NewContext(recorder, nil)

	require.Equal(t, "", recorder.Body.String())

	// Act
	context.WriteBody([]byte("test body"))

	// Assert
	require.Equal(t, "test body", recorder.Body.String())
}

func TestContext_SetHeader(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name string

		initialHeader http.Header

		headerKeyToAdd   string
		headerValueToAdd string

		expectedHeader http.Header
	}

	tests := []TestCase{
		{
			name:             "GivenKeyIsMissing_ShouldSetValue",
			initialHeader:    make(http.Header),
			headerKeyToAdd:   "Content-Type",
			headerValueToAdd: "application/json",
			expectedHeader: http.Header{
				"Content-Type": []string{"application/json"},
			},
		},
		{
			name: "GivenKeyIsPresent_ValueContainsASingleElement_ShouldReplaceValueElement",
			initialHeader: http.Header{
				"Key": []string{"Value1"},
			},
			headerKeyToAdd:   "Key",
			headerValueToAdd: "Value",
			expectedHeader: http.Header{
				"Key": []string{"Value"},
			},
		},
		{
			name: "GivenKeyIsPresent_ValueContainsMultipleElements_ShouldReplaceValueToSingleElementSlice",
			initialHeader: http.Header{
				"Key": []string{"Value1", "Value2"},
			},
			headerKeyToAdd:   "Key",
			headerValueToAdd: "Value",
			expectedHeader: http.Header{
				"Key": []string{"Value"},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			recorder := httptest.NewRecorder()
			header := recorder.Header()

			for key, values := range test.initialHeader {
				header[key] = values
			}

			context := lit.NewContext(recorder, nil)

			// Act
			context.SetHeader(test.headerKeyToAdd, test.headerValueToAdd)

			// Assert
			require.Equal(t, test.expectedHeader, header)
		})
	}
}
