package bind_test

import (
	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/bind"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestQuery_WhenTypeParameterIsNotStruct_ShouldPanic(t *testing.T) {
	t.Parallel()

	request := lit.NewRequest(
		httptest.NewRequest(http.MethodGet, "/", nil),
		nil,
	)

	require.PanicsWithValue(t, "T must be a struct type",
		func() { _, _ = bind.Query[int](request) })
}

func TestQuery(t *testing.T) {
	t.Parallel()

	type targetStruct struct {
		UserID    int         `query:"user_id"`
		BookID    string      `query:"book_id"`
		TimeRange []time.Time `query:"time_range"`
		Store     string
	}

	tests := []struct {
		description    string
		parameters     url.Values
		expectedResult any
		expectedError  string
	}{
		{
			description: "WhenQueryParametersMatchBindingTarget_ShouldBind",
			parameters: map[string][]string{
				"user_id":    {"123"},
				"book_id":    {"Book Name"},
				"time_range": {"2023-10-22T00:00:00Z", "2023-10-25T23:59:00Z"},
			},
			expectedResult: targetStruct{
				UserID: 123,
				BookID: "Book Name",
				TimeRange: []time.Time{
					time.Date(2023, 10, 22, 00, 00, 00, 00, time.UTC),
					time.Date(2023, 10, 25, 23, 59, 00, 00, time.UTC),
				},
			},
			expectedError: "",
		},
		{
			description: "WhenQueryParametersAreMissingInRequest_ShouldIgnoreThem",
			parameters: map[string][]string{
				"user_id": {"123"},
			},
			expectedResult: targetStruct{
				UserID: 123,
				BookID: "",
			},
			expectedError: "",
		},
		{
			description: "WhenQueryParametersDoNotMatchBindingTarget_ShouldReturnError",
			parameters: map[string][]string{
				"user_id": {"123a"},
			},
			expectedResult: targetStruct{
				UserID: 0,
				BookID: "",
			},
			expectedError: "user_id: 123a is not a valid int: invalid syntax",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.URL.RawQuery = test.parameters.Encode()

			request := lit.NewRequest(r, nil)

			// Act
			result, err := bind.Query[targetStruct](request)

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
