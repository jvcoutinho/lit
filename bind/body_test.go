package bind_test

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/bind"
	"github.com/stretchr/testify/require"
)

func TestBody(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description    string
		body           string
		function       func(r *lit.Request) (any, error)
		contentType    string
		expectedResult any
		expectedError  string
		shouldPanic    bool
	}{
		{
			description: "WhenTypeParameterIsNotStruct_ShouldPanic",
			function: func(r *lit.Request) (any, error) {
				return bind.Body[int](r)
			},
			expectedError: "T must be a struct type",
			shouldPanic:   true,
		},
		{
			description: "WhenContentTypeIsUnsupported_ShouldReturnError",
			contentType: "text/plain",
			function: func(r *lit.Request) (any, error) {
				return bind.Body[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "unsupported Content-Type",
		},
		{
			description: "WhenContentTypeIsJSON_AndItIsValid_ShouldBind",
			body: `
{
    "string": "hi",
	"pointer": 10,
    "uint": 10,
    "uint8": 10,
    "uint16": 10,
    "uint32": 10,
    "uint64": 10,
    "int": 10,
    "int8": 10,
    "int16": 10,
    "int32": 10,
    "int64": 10,
    "float32": 10.5,
    "float64": 10.5,
    "bool": true,
    "time": "2023-10-22T00:00:00Z",
    "slice": [
        2,
        3
    ],
    "array": [
        2,
        3
    ]
}`,
			function: func(r *lit.Request) (any, error) {
				return bind.Body[bindableFields](r)
			},
			expectedResult: bindableFields{
				String:  "hi",
				Pointer: pointerOf(10),
				Uint:    10,
				Uint8:   10,
				Uint16:  10,
				Uint32:  10,
				Uint64:  10,
				Int:     10,
				Int8:    10,
				Int16:   10,
				Int32:   10,
				Int64:   10,
				Float32: 10.5,
				Float64: 10.5,
				Bool:    true,
				Time:    time.Date(2023, 10, 22, 0, 0, 0, 0, time.UTC),
				Slice:   []int{2, 3},
				Array:   [2]int{2, 3},
			},
		},
		{
			description: "WhenContentTypeIsJSON_AndItIsInvalid_ShouldReturnError",
			body: `
{
    "uint": "10a"
}`,
			function: func(r *lit.Request) (any, error) {
				return bind.Body[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "uint: string is not a valid uint",
		},
		{
			description: "WhenBodyIsEmpty_ShouldReturnDefaultValue",
			body:        "",
			function: func(r *lit.Request) (any, error) {
				return bind.Body[bindableFields](r)
			},
			expectedResult: bindableFields{},
		},
		{
			description: "WhenContentTypeIsYAML_AndItIsValid_ShouldBind",
			body: `
string: hi
pointer: 10
uint: 10
uint8: 10
uint16: 10
uint32: 10
uint64: 10
int: 10
int8: 10
int16: 10
int32: 10
int64: 10
float32: 10.5
float64: 10.5
bool: true
time: '2023-10-22T00:00:00Z'
slice:
- 2
- 3
array:
- 2
- 3
`,
			contentType: "application/x-yaml",
			function: func(r *lit.Request) (any, error) {
				return bind.Body[bindableFields](r)
			},
			expectedResult: bindableFields{
				String:  "hi",
				Pointer: pointerOf(10),
				Uint:    10,
				Uint8:   10,
				Uint16:  10,
				Uint32:  10,
				Uint64:  10,
				Int:     10,
				Int8:    10,
				Int16:   10,
				Int32:   10,
				Int64:   10,
				Float32: 10.5,
				Float64: 10.5,
				Bool:    true,
				Time:    time.Date(2023, 10, 22, 0, 0, 0, 0, time.UTC),
				Slice:   []int{2, 3},
				Array:   [2]int{2, 3},
			},
		},
		{
			description: "WhenContentTypeIsYAML_AndItIsInvalid_ShouldReturnError",
			body:        `uint: 10a`,
			contentType: "application/x-yaml",
			function: func(r *lit.Request) (any, error) {
				return bind.Body[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `10a` into uint",
		},
		{
			description: "WhenContentTypeIsXML_AndItIsValid_ShouldBind",
			body: `
<?xml version="1.0" encoding="UTF-8" ?>
<root>
  <String>hi</String>
  <Pointer>10</Pointer>
  <Uint>10</Uint>
  <Uint8>10</Uint8>
  <Uint16>10</Uint16>
  <Uint32>10</Uint32>
  <Uint64>10</Uint64>
  <Int>10</Int>
  <Int8>10</Int8>
  <Int16>10</Int16>
  <Int32>10</Int32>
  <Int64>10</Int64>
  <Float32>10.5</Float32>
  <Float64>10.5</Float64>
  <Bool>true</Bool>
  <Time>2023-10-22T00:00:00Z</Time>
  <Slice>2</Slice>
  <Slice>3</Slice>
</root>
`,
			contentType: "application/xml",
			function: func(r *lit.Request) (any, error) {
				return bind.Body[bindableFields](r)
			},
			expectedResult: bindableFields{
				String:  "hi",
				Pointer: pointerOf(10),
				Uint:    10,
				Uint8:   10,
				Uint16:  10,
				Uint32:  10,
				Uint64:  10,
				Int:     10,
				Int8:    10,
				Int16:   10,
				Int32:   10,
				Int64:   10,
				Float32: 10.5,
				Float64: 10.5,
				Bool:    true,
				Time:    time.Date(2023, 10, 22, 0, 0, 0, 0, time.UTC),
				Slice:   []int{2, 3},
			},
		},
		{
			description: "WhenContentTypeIsXML_AndItIsInvalid_ShouldReturnError",
			body: `
<?xml version="1.0" encoding="UTF-8" ?>
<root>
  <Uint>10a</Uint>
</root>
`,
			contentType: "application/xml",
			function: func(r *lit.Request) (any, error) {
				return bind.Body[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  `strconv.ParseUint: parsing "10a": invalid syntax`,
		},
		{
			description: "WhenContentTypeIsForm_AndItIsValid_ShouldBind",
			body: "string=hi&pointer=10&uint=10&uint8=10&uint16=10&uint32=10&uint64=10&int=10&int8=10&int16=10" +
				"&int32=10&int64=10&float32=10.5&float64=10.5&complex64=10.5&complex128=10.5&bool=true" +
				"&time=2023-10-22T00:00:00Z&slice=2&slice=3&array=2&array=3",
			contentType: "application/x-www-form-urlencoded",
			function: func(r *lit.Request) (any, error) {
				return bind.Body[bindableFields](r)
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
			description: "WhenContentTypeIsForm_AndItIsInvalid_ShouldReturnError",
			body:        "uint=10a",
			contentType: "application/x-www-form-urlencoded",
			function: func(r *lit.Request) (any, error) {
				return bind.Body[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "uint: 10a is not a valid uint: invalid syntax",
		},
		{
			description: "WhenContentTypeIsForm_AndItIsMalformed_ShouldReturnError",
			body:        "&key;&",
			contentType: "application/x-www-form-urlencoded",
			function: func(r *lit.Request) (any, error) {
				return bind.Body[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "invalid semicolon separator in query",
		},
		{
			description: "WhenContentTypeIsMultipartForm_AndItIsValid_ShouldBind",
			body: `
--BOUNDARY
Content-Disposition: form-data; name="string"

hi
--BOUNDARY
Content-Disposition: form-data; name="pointer"

10
--BOUNDARY
Content-Disposition: form-data; name="uint"

10
--BOUNDARY
Content-Disposition: form-data; name="uint8"

10
--BOUNDARY
Content-Disposition: form-data; name="uint16"

10
--BOUNDARY
Content-Disposition: form-data; name="uint32"

10
--BOUNDARY
Content-Disposition: form-data; name="uint64"

10
--BOUNDARY
Content-Disposition: form-data; name="int"

10
--BOUNDARY
Content-Disposition: form-data; name="int8"

10
--BOUNDARY
Content-Disposition: form-data; name="int16"

10
--BOUNDARY
Content-Disposition: form-data; name="int32"

10
--BOUNDARY
Content-Disposition: form-data; name="int64"

10
--BOUNDARY
Content-Disposition: form-data; name="float32"

10.5
--BOUNDARY
Content-Disposition: form-data; name="float64"

10.5
--BOUNDARY
Content-Disposition: form-data; name="complex64"

10.5
--BOUNDARY
Content-Disposition: form-data; name="complex128"

10.5
--BOUNDARY
Content-Disposition: form-data; name="bool"

true
--BOUNDARY
Content-Disposition: form-data; name="time"

2023-10-22T00:00:00Z
--BOUNDARY
Content-Disposition: form-data; name="slice"

2
--BOUNDARY
Content-Disposition: form-data; name="slice"

3
--BOUNDARY
Content-Disposition: form-data; name="array"

2
--BOUNDARY
Content-Disposition: form-data; name="array"

3
--BOUNDARY--
`,
			contentType: "multipart/form-data; boundary=BOUNDARY",
			function: func(r *lit.Request) (any, error) {
				return bind.Body[bindableFields](r)
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
			description: "WhenContentTypeIsMultipartForm_AndItHasFile_ShouldBind",
			body: `
--BOUNDARY
Content-Disposition: form-data; name="file"; filename="file.txt"
Content-Type: text/plain

content of a file
--BOUNDARY
Content-Disposition: form-data; name="files"; filename="file1.txt"
Content-Type: text/plain

contents of file 1
--BOUNDARY
Content-Disposition: form-data; name="files"; filename="file2.txt"
Content-Type: text/plain

contents of file 2
--BOUNDARY--`,
			contentType: "multipart/form-data; boundary=BOUNDARY",
			function: func(r *lit.Request) (any, error) {
				return bind.Body[bindableFields](r)
			},
			expectedResult: bindableFields{
				File: &multipart.FileHeader{
					Filename: "file.txt",
					Header: textproto.MIMEHeader{
						"Content-Disposition": {`form-data; name="file"; filename="file.txt"`},
						"Content-Type":        {"text/plain"},
					},
					Size: 17,
				},
				Files: []*multipart.FileHeader{
					{
						Filename: "file1.txt",
						Header: textproto.MIMEHeader{
							"Content-Disposition": {`form-data; name="files"; filename="file1.txt"`},
							"Content-Type":        {"text/plain"},
						},
						Size: 18,
					},
					{
						Filename: "file2.txt",
						Header: textproto.MIMEHeader{
							"Content-Disposition": {`form-data; name="files"; filename="file2.txt"`},
							"Content-Type":        {"text/plain"},
						},
						Size: 18,
					},
				},
			},
		},
		{
			description: "WhenContentTypeIsMultipartForm_AndItHasFile_ButTypeIsIncorrect_ShouldReturnError",
			body: `
--BOUNDARY
Content-Disposition: form-data; name="file"; filename="file.txt"
Content-Type: text/plain

content of a file
--BOUNDARY--`,
			contentType: "multipart/form-data; boundary=BOUNDARY",
			function: func(r *lit.Request) (any, error) {
				return bind.Body[unbindableField](r)
			},
			expectedResult: unbindableField{},
			expectedError:  "file: [file.txt] is not a valid int",
		},
		{
			description: "WhenContentTypeIsMultipartForm_AndItIsInvalid_ShouldReturnError",
			body:        "--BOUNDARY\nContent-Disposition: form-data; name=\"uint\"\n\n10a\n--BOUNDARY--",
			contentType: "multipart/form-data; boundary=BOUNDARY",
			function: func(r *lit.Request) (any, error) {
				return bind.Body[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  "uint: 10a is not a valid uint: invalid syntax",
		},
		{
			description: "WhenContentTypeIsMultipartForm_AndItIsMalformed_ShouldReturnError",
			body:        "--BOUNDARY\n&key;&\n--BOUNDARY--",
			contentType: "multipart/form-data; boundary=BOUNDARY",
			function: func(r *lit.Request) (any, error) {
				return bind.Body[bindableFields](r)
			},
			expectedResult: bindableFields{},
			expectedError:  `malformed MIME header: missing colon: "&key;&"`,
		},
		{
			description: "WhenTypeParameterIsValidatableWithValueReceiver_ShouldNotValidate",
			body:        `{"string": "string"}`,
			function: func(r *lit.Request) (any, error) {
				return bind.Body[nonPointerReceiverValidatableFields](r)
			},
			expectedResult: nonPointerReceiverValidatableFields{String: "string"},
		},
		{
			description: "WhenTypeParameterIsValidatableWithPointerReceiver_ShouldValidate",
			body:        `{"string": "string"}`,
			function: func(r *lit.Request) (any, error) {
				return bind.Body[pointerReceiverValidatableFields](r)
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
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.body))
			request.Header.Add("Content-Type", test.contentType)

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
			require.True(t,
				cmp.Equal(test.expectedResult, result, cmpopts.IgnoreUnexported(multipart.FileHeader{}, lit.Request{})),
				cmp.Diff(test.expectedResult, result, cmpopts.IgnoreUnexported(multipart.FileHeader{}, lit.Request{})),
			)
		})
	}
}

func ExampleBody() {
	req := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(`
		{"name": "Percy Jackson", "publishYear": 2009}
	`))

	r := lit.NewRequest(req)

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

func ExampleBody_form() {
	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = "publishYear=2009&name=Percy%20Jackson"

	r := lit.NewRequest(req)

	type RequestBody struct {
		Name        string `form:"name"`
		PublishYear int    `form:"publishYear"`
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
