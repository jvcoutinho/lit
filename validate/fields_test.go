package validate_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/bind"

	"github.com/jvcoutinho/lit/validate"

	"github.com/stretchr/testify/require"
)

func TestFields(t *testing.T) {
	t.Parallel()

	validatableFields := validatableFields{
		String:  "string",
		Int:     30,
		Pointer: pointerOf(5),
	}

	tests := []struct {
		description   string
		function      func() error
		expectedError string
		shouldPanic   bool
	}{
		{
			description: "WhenThereAreNoValidations_ShouldNotReturnError",
			function: func() error {
				return validate.Fields(&validatableFields)
			},
		},
		{
			description: "WhenThereAreNoViolationsShouldNotReturnError",
			function: func() error {
				return validate.Fields(&validatableFields,
					validate.Field{
						Valid: validatableFields.Int <= 60,
					},
					validate.Field{
						Valid: validatableFields.Int >= 30,
					},
				)
			},
		},
		{
			description: "WhenThereAreViolations_AndTypeParameterIsNotStruct_ShouldPanic",
			function: func() error {
				value := 2

				return validate.Fields(&value,
					validate.Field{
						Valid: validatableFields.Int > 60,
					},
				)
			},
			expectedError: "T must be a struct type",
			shouldPanic:   true,
		},
		{
			description: "WhenThereAreViolations_AndValidationFieldsAreNotPointers_ShouldPanic",
			function: func() error {
				value := 2

				return validate.Fields(&validatableFields,
					validate.Field{
						Valid:  validatableFields.Int > 60,
						Fields: []any{value},
					},
				)
			},
			expectedError: "argument 0 should be a pointer to a field of validate_test.validatableFields",
			shouldPanic:   true,
		},
		{
			description: "WhenThereAreViolations_AndValidationFieldsAreNotFieldPointers_ShouldPanic",
			function: func() error {
				value := 2

				return validate.Fields(&validatableFields,
					validate.Field{
						Valid:  validatableFields.Int > 60,
						Fields: []any{&value},
					},
				)
			},
			expectedError: "argument 0 should be a pointer to a field of validate_test.validatableFields",
			shouldPanic:   true,
		},
		{
			description: "WhenThereAreViolations_AndTagIsPresent_ShouldBuildErrorMessageWithTagValue",
			function: func() error {
				return validate.Fields(&validatableFields,
					validate.Field{
						Valid:   validatableFields.Int > 60,
						Message: "{0} should be greater than 60",
						Fields:  []any{&validatableFields.Int},
					},
				)
			},
			expectedError: "int should be greater than 60",
		},
		{
			description: "WhenThereAreViolations_AndTagIsEmpty_ShouldBuildErrorMessageWithFieldName",
			function: func() error {
				return validate.Fields(&validatableFields,
					validate.Field{
						Valid:   validatableFields.Bool,
						Message: "{0} should be true",
						Fields:  []any{&validatableFields.Bool},
					},
				)
			},
			expectedError: "Bool should be true",
		},
		{
			description: "WhenThereAreViolations_AndTagIsNotPresent_ShouldBuildErrorMessageWithFieldName",
			function: func() error {
				return validate.Fields(&validatableFields,
					validate.Field{
						Valid:   len(validatableFields.String) > 20,
						Message: "{0} should have a length greater than 20",
						Fields:  []any{&validatableFields.String},
					},
				)
			},
			expectedError: "String should have a length greater than 20",
		},
		{
			description: "WhenThereAreViolations_ShouldBuildACombinedErrorMessage",
			function: func() error {
				return validate.Fields(&validatableFields,
					validate.Field{
						Valid:   len(validatableFields.String) > 20,
						Message: "{0} should have a length greater than 20",
						Fields:  []any{&validatableFields.String},
					},
					validate.Field{
						Valid:   strings.HasPrefix(validatableFields.String, "a"),
						Message: `{0} should start with prefix "a"`,
						Fields:  []any{&validatableFields.String},
					},
				)
			},
			expectedError: `String should have a length greater than 20; String should start with prefix "a"`,
		},
		{
			description: "WhenThereArePointerViolations_ShouldValidateItsValue",
			function: func() error {
				return validate.Fields(&validatableFields,
					validate.Field{
						Valid:   *validatableFields.Pointer == 2,
						Message: "{0} should be equal to 2",
						Fields:  []any{validatableFields.Pointer},
					},
				)
			},
			expectedError: `pointer should be equal to 2`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			if test.shouldPanic {
				require.PanicsWithValue(t, test.expectedError, func() {
					_ = test.function()
				})

				return
			}

			err := test.function()

			// Assert
			errMessage := ""
			if err != nil {
				errMessage = err.Error()
			}

			require.Equal(t, test.expectedError, errMessage)
		})
	}
}

func ExampleFields() {
	req := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(`
		{"name": "Percy Jackson", "publishYear": 2007}
	`))

	r := lit.NewRequest(req)

	type RequestBody struct {
		Name        string `json:"name" validate:"name"`
		PublishYear int    `json:"publishYear" validate:"publishYear"`
	}

	body, _ := bind.Body[RequestBody](r)

	err := validate.Fields(&body,
		validate.Greater(&body.PublishYear, 2009),
		validate.Less(&body.PublishYear, 2020),
		validate.HasPrefix(&body.Name, "A"),
	)

	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// publishYear should be greater than 2009; name should start with "A"
}
