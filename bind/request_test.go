package bind_test

import (
	"bytes"
	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/bind"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestRequest_WhenTargetTypeIsNotStruct_ShouldPanic(t *testing.T) {
	t.Parallel()

	request := lit.NewRequest(nil, nil)

	require.PanicsWithValue(t, "int is not a struct type", func() {
		_, _ = bind.Request[int](request)
	})
}

func TestRequest(t *testing.T) {
	t.Parallel()

	type targetStruct struct {
		Authorization string `header:"Authorization"`
		UserID        int    `uri:"user_id"`
		Name          string `json:"name"`
		PublishYear   int    `json:"publish_year"`
	}

	tests := []struct {
		description     string
		header          http.Header
		body            string
		queryParameters url.Values
		urlParameters   map[string]string
		expectedResult  any
		expectedError   string
	}{
		{
			description: "WhenArgumentsMatchBindingTarget_ShouldBind",
			header: map[string][]string{
				"Content-Type":   {"application/json"},
				"Authorization":  {"Bearer token"},
				"Content-Length": {"32"},
			},
			body:            `{"name": "Alive and Well", "publish_year": 2008}`,
			queryParameters: nil,
			urlParameters: map[string]string{
				"user_id": "123",
			},
			expectedResult: targetStruct{
				Authorization: "Bearer token",
				UserID:        123,
				Name:          "Alive and Well",
				PublishYear:   2008,
			},
			expectedError: "",
		},
		{
			description: "WhenContentIsMissingInRequest_ShouldIgnoreIt",
			header: map[string][]string{
				"Content-Type": {"application/json"},
			},
			expectedResult: targetStruct{},
			expectedError:  "",
		},
		{
			description:    "WhenContentDoesNotMatchBindingTarget_ShouldReturnError",
			urlParameters:  map[string]string{"user_id": "123a"},
			expectedResult: targetStruct{},
			expectedError:  "user_id: 123a is not a valid int: invalid syntax",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			r := httptest.NewRequest(http.MethodGet, "/123", bytes.NewReader([]byte(test.body)))
			r.Header = test.header
			r.URL.RawQuery = test.queryParameters.Encode()

			request := lit.NewRequest(r, test.urlParameters)

			// Act
			result, err := bind.Request[targetStruct](request)

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
