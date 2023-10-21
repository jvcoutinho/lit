package bind_test

import (
	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/bind"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestURLParameters_WhenTargetTypeIsNotStruct_ShouldPanic(t *testing.T) {
	t.Parallel()

	request := lit.NewRequest(nil, nil)

	require.PanicsWithValue(t, "int is not a struct type", func() {
		_, _ = bind.URLParameters[int](request)
	})
}

func TestURLParameters(t *testing.T) {
	t.Parallel()

	type targetStruct struct {
		UserID int    `argument:"user_id"`
		BookID string `argument:"book_id"`
		Store  string
	}

	tests := []struct {
		description    string
		arguments      map[string]string
		expectedResult any
		expectedError  string
	}{
		{
			description: "WhenArgumentsMatchBindingTarget_ShouldBind",
			arguments: map[string]string{
				"user_id": "123",
				"book_id": "Book Name",
			},
			expectedResult: targetStruct{
				UserID: 123,
				BookID: "Book Name",
			},
			expectedError: "",
		},
		{
			description: "WhenArgumentsAreMissing_ShouldIgnoreThem",
			arguments: map[string]string{
				"user_id": "123",
			},
			expectedResult: targetStruct{
				UserID: 123,
				BookID: "",
			},
			expectedError: "",
		},
		{
			description: "WhenArgumentsDoNotMatchBindingTarget_ShouldReturnError",
			arguments: map[string]string{
				"user_id": "123a",
			},
			expectedResult: targetStruct{
				UserID: 0,
				BookID: "",
			},
			expectedError: "binding into user_id: 123a is not a valid int",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			request := lit.NewRequest(nil, test.arguments)

			// Act
			result, err := bind.URLParameters[targetStruct](request)

			// Assert
			require.Equal(t, test.expectedResult, result)

			if test.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, test.expectedError)
			}
		})
	}
}
