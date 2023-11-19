package bind_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/bind"
	"github.com/stretchr/testify/require"
)

func TestBody(t *testing.T) {
	t.Parallel()

	type result struct {
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
			description: "Valid JSON",
			body: `
				{
					"name": "John",
					"age": 27,
					"is_married": true
				}
			`,
			contentType: "application/json",
			expectedResult: result{
				Name:      "John",
				Age:       27,
				IsMarried: true,
			},
			expectedError: "",
		},
		{
			description: "Invalid JSON",
			body: `
				{
					"name": John",
					"age": 27,
					"is_married": true
				}
			`,
			contentType:    "application/json",
			expectedResult: result{},
			expectedError:  "invalid character 'J' looking for beginning of value",
		},
		{
			description: "Valid YAML 1",
			body: `
name: John
age: 27
is_married: true`,
			contentType: "application/x-yaml",
			expectedResult: result{
				Name:      "John",
				Age:       27,
				IsMarried: true,
			},
			expectedError: "",
		},
		{
			description: "Valid YAML 2",
			body: `
name: John
age: 27
is_married: true`,
			contentType: "text/yaml",
			expectedResult: result{
				Name:      "John",
				Age:       27,
				IsMarried: true,
			},
			expectedError: "",
		},
		{
			description: "Invalid YAML 1",
			body: `
name: John
 age: 27
is_married: true`,
			contentType:    "application/x-yaml",
			expectedResult: result{},
			expectedError:  "yaml: line 3: mapping values are not allowed in this context",
		},
		{
			description: "Invalid YAML 2",
			body: `
name: John
 age: 27
is_married: true`,
			contentType:    "text/yaml",
			expectedResult: result{},
			expectedError:  "yaml: line 3: mapping values are not allowed in this context",
		},
		{
			description: "Valid XML 1",
			body: `
<?xml version="1.0" encoding="UTF-8" ?>
<root>
  <Name>John</Name>
  <Age>27</Age>
  <IsMarried>true</IsMarried>
</root>
			`,
			contentType: "application/xml",
			expectedResult: result{
				Name:      "John",
				Age:       27,
				IsMarried: true,
			},
			expectedError: "",
		},
		{
			description: "Valid XML 2",
			body: `
<?xml version="1.0" encoding="UTF-8" ?>
<root>
  <Name>John</Name>
  <Age>27</Age>
  <IsMarried>true</IsMarried>
</root>
			`,
			contentType: "text/xml",
			expectedResult: result{
				Name:      "John",
				Age:       27,
				IsMarried: true,
			},
			expectedError: "",
		},
		{
			description: "Invalid XML 1",
			body: `
<?xml version="1.0" encoding="UTF-8" ?>
<roots>
  <Name>John</Name>
  <Age>27</Age>
  <IsMarried>true</IsMarried>
</root>
			`,
			contentType: "application/xml",
			expectedResult: result{
				Name:      "John",
				Age:       27,
				IsMarried: true,
			},
			expectedError: "XML syntax error on line 7: element <roots> closed by </root>",
		},
		{
			description: "Invalid XML 2",
			body: `
<?xml version="1.0" encoding="UTF-8" ?>
<roots>
  <Name>John</Name>
  <Age>27</Age>
  <IsMarried>true</IsMarried>
</root>
			`,
			contentType: "text/xml",
			expectedResult: result{
				Name:      "John",
				Age:       27,
				IsMarried: true,
			},
			expectedError: "XML syntax error on line 7: element <roots> closed by </root>",
		},
		{
			description: "MissingContentType_ValidJSON",
			body: `
				{
					"name": "John",
					"age": 27,
					"is_married": true
				}
			`,
			contentType: "",
			expectedResult: result{
				Name:      "John",
				Age:       27,
				IsMarried: true,
			},
			expectedError: "",
		},
		{
			description: "MissingContentType_InvalidJSON",
			body: `
name: John
age: 27
is_married: true`,
			contentType:    "",
			expectedResult: result{},
			expectedError:  "invalid character 'a' in literal null (expecting 'u')",
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
			result, err := bind.Body[result](request)

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

func ExampleBody() {
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBufferString(`
		{"name": "Percy Jackson", "publishYear": 2009}
	`))

	r := lit.NewRequest(req, nil)

	type RequestBody struct {
		Name        string `json:"name"`
		PublishYear int    `json:"publishYear"`
	}

	body, err := bind.Body[RequestBody](r)
	if err == nil {
		fmt.Println(body.Name)
		fmt.Println(body.PublishYear)
	}

	// Output:
	// Percy Jackson
	// 2009
}
