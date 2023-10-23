package bind_test

import (
	"bytes"
	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/bind"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBody(t *testing.T) {
	t.Parallel()

	type targetStruct struct {
		Name      string `json:"name" yaml:"name" xml:"Name"`
		Age       int    `json:"age" yaml:"age" xml:"Age"`
		IsMarried bool   `json:"is_married" yaml:"is_married" xml:"IsMarried"`
	}

	tests := []struct {
		description    string
		body           string
		contentType    string
		expectedResult any
		expectedError  string
	}{
		{
			description: "JSON_WhenArgumentsMatchBindingTarget_ShouldBind",
			body: `
				{
					"name": "John",
					"age": 27,
					"is_married": true
				}
			`,
			contentType: "application/json",
			expectedResult: targetStruct{
				Name:      "John",
				Age:       27,
				IsMarried: true,
			},
			expectedError: "",
		},
		{
			description: "YAML_WhenArgumentsMatchBindingTarget_ShouldBind",
			body: `
name: John
age: 27
is_married: true`,
			contentType: "application/x-yaml",
			expectedResult: targetStruct{
				Name:      "John",
				Age:       27,
				IsMarried: true,
			},
			expectedError: "",
		},
		{
			description: "XML_WhenArgumentsMatchBindingTarget_ShouldBind",
			body: `
<?xml version="1.0" encoding="UTF-8" ?>
<root>
  <Name>John</Name>
  <Age>27</Age>
  <IsMarried>true</IsMarried>
</root>
			`,
			contentType: "application/xml",
			expectedResult: targetStruct{
				Name:      "John",
				Age:       27,
				IsMarried: true,
			},
			expectedError: "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(test.body)))
			r.Header.Add("Content-Type", test.contentType)

			request := lit.NewRequest(r, nil)

			// Act
			result, err := bind.Body[targetStruct](request)

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
