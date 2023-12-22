package render_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jvcoutinho/lit/render"
	"github.com/stretchr/testify/require"
)

func TestJSONResponse(t *testing.T) {
	t.Parallel()

	type jsonObject struct {
		Key1 string `json:"key_1"`
		Key2 int    `json:"key_2"`
		Key3 bool   `json:"key_3"`
	}

	tests := []struct {
		description        string
		response           render.JSONResponse
		expectedStatusCode int
		expectedBody       string
		expectedHeader     http.Header
	}{
		{
			description: "WhenBodyIsObject_ShouldMarshalAsIs",
			response: render.JSON(203, jsonObject{
				Key1: "key1",
				Key2: 100,
				Key3: true,
			}),
			expectedStatusCode: 203,
			expectedBody:       `{"key_1":"key1","key_2":100,"key_3":true}`,
			expectedHeader:     http.Header{"Content-Type": []string{"application/json"}},
		},
		{
			description:        "WhenBodyIsString_ShouldMarshalMessageResponse",
			response:           render.JSON(203, "test string"),
			expectedStatusCode: 203,
			expectedBody:       `{"message":"test string"}`,
			expectedHeader:     http.Header{"Content-Type": []string{"application/json"}},
		},
		{
			description:        "WhenBodyIsError_ShouldMarshalMessageResponse",
			response:           render.JSON(203, errors.New("test error")),
			expectedStatusCode: 203,
			expectedBody:       `{"message":"test error"}`,
			expectedHeader:     http.Header{"Content-Type": []string{"application/json"}},
		},
		{
			description:        "WhenBodyIsNil_ShouldPrintNoBody",
			response:           render.JSON(203, nil),
			expectedStatusCode: 203,
			expectedBody:       "",
			expectedHeader:     http.Header{},
		},
		{
			description:        "WhenBodyCanNotBeMarshalled_ShouldRespondInternalServerError",
			response:           render.JSON(203, func() {}),
			expectedStatusCode: 500,
			expectedBody:       "json: unsupported type: func()\n",
			expectedHeader: http.Header{
				"Content-Type":           []string{"text/plain; charset=utf-8"},
				"X-Content-Type-Options": []string{"nosniff"},
			},
		},
		{
			description:        "ShouldSetResponseHeader",
			response:           render.JSON(203, "body").WithHeader("Key", "Value"),
			expectedStatusCode: 203,
			expectedBody:       `{"message":"body"}`,
			expectedHeader: http.Header{
				"Content-Type": []string{"application/json"},
				"Key":          []string{"Value"},
			},
		},
		{
			description:        "OK_ShouldMarshalJSONResponse",
			response:           render.OK("body"),
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message":"body"}`,
			expectedHeader:     http.Header{"Content-Type": []string{"application/json"}},
		},
		{
			description:        "Created_ShouldMarshalJSONResponse",
			response:           render.Created("body", "/location-resource"),
			expectedStatusCode: http.StatusCreated,
			expectedBody:       `{"message":"body"}`,
			expectedHeader: http.Header{
				"Content-Type": []string{"application/json"},
				"Location":     []string{"/location-resource"},
			},
		},
		{
			description:        "Accepted_ShouldMarshalJSONResponse",
			response:           render.Accepted("body"),
			expectedStatusCode: http.StatusAccepted,
			expectedBody:       `{"message":"body"}`,
			expectedHeader:     http.Header{"Content-Type": []string{"application/json"}},
		},
		{
			description:        "BadRequest_ShouldMarshalJSONResponse",
			response:           render.BadRequest("body"),
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"message":"body"}`,
			expectedHeader:     http.Header{"Content-Type": []string{"application/json"}},
		},
		{
			description:        "Unauthorized_ShouldMarshalJSONResponse",
			response:           render.Unauthorized("body"),
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"message":"body"}`,
			expectedHeader:     http.Header{"Content-Type": []string{"application/json"}},
		},
		{
			description:        "Forbidden_ShouldMarshalJSONResponse",
			response:           render.Forbidden("body"),
			expectedStatusCode: http.StatusForbidden,
			expectedBody:       `{"message":"body"}`,
			expectedHeader:     http.Header{"Content-Type": []string{"application/json"}},
		},
		{
			description:        "NotFound_ShouldMarshalJSONResponse",
			response:           render.NotFound("body"),
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       `{"message":"body"}`,
			expectedHeader:     http.Header{"Content-Type": []string{"application/json"}},
		},
		{
			description:        "Conflict_ShouldMarshalJSONResponse",
			response:           render.Conflict("body"),
			expectedStatusCode: http.StatusConflict,
			expectedBody:       `{"message":"body"}`,
			expectedHeader:     http.Header{"Content-Type": []string{"application/json"}},
		},
		{
			description:        "UnprocessableContent_ShouldMarshalJSONResponse",
			response:           render.UnprocessableContent("body"),
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedBody:       `{"message":"body"}`,
			expectedHeader:     http.Header{"Content-Type": []string{"application/json"}},
		},
		{
			description:        "InternalServerError_ShouldMarshalJSONResponse",
			response:           render.InternalServerError("body"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"message":"body"}`,
			expectedHeader:     http.Header{"Content-Type": []string{"application/json"}},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			var (
				writer   = httptest.NewRecorder()
				response = test.response
			)

			// Act
			response.Write(writer)

			// Assert
			require.Equal(t, test.expectedStatusCode, writer.Code)
			require.Equal(t, test.expectedHeader, writer.Header())
			require.Equal(t, test.expectedBody, writer.Body.String())
		})
	}
}
