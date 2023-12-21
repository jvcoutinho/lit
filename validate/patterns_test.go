package validate_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jvcoutinho/lit/validate"
	"github.com/stretchr/testify/require"
)

func TestUUID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *string
		expected    bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			expected:    false,
		},
		{
			description: "IsUUIDv1",
			target:      pointerOf("447aa620-a03a-11ee-9c11-8d9995a85f9f"),
			expected:    true,
		},
		{
			description: "IsUUIDv4",
			target:      pointerOf("212c2105-22f6-4979-bf3c-536939392aeb"),
			expected:    true,
		},
		{
			description: "IsUUIDv5",
			target:      pointerOf("4e6ced17-7791-5536-8ebe-c992f82d2826"),
			expected:    true,
		},
		{
			description: "IsNotUUID",
			target:      pointerOf("not-uuid"),
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.UUID(test.target)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: "{0} is not a valid UUID",
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}

func TestEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *string
		expected    bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			expected:    false,
		},
		{
			description: "IsEmail",
			target:      pointerOf("username@domain.com"),
			expected:    true,
		},
		{
			description: "DoesNotHaveUsername",
			target:      pointerOf("@domain.com"),
			expected:    false,
		},
		{
			description: "DoesNotHaveDomain",
			target:      pointerOf("username@"),
			expected:    false,
		},
		{
			description: "DoesNotHave@",
			target:      pointerOf("usernamedomain.com"),
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.Email(test.target)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: "{0} is not a valid e-mail",
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}

func TestIPAddress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *string
		expected    bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			expected:    false,
		},
		{
			description: "IsIPV4",
			target:      pointerOf("186.15.46.163"),
			expected:    true,
		},
		{
			description: "IsIPV6",
			target:      pointerOf("007e:89eb:9dc0:73cd:326f:503e:38d5:697e"),
			expected:    true,
		},
		{
			description: "IsNotIP",
			target:      pointerOf("string"),
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.IPAddress(test.target)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: "{0} is not a valid IP address",
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}

func TestIPv4Address(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *string
		expected    bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			expected:    false,
		},
		{
			description: "IsIPV4",
			target:      pointerOf("186.15.46.163"),
			expected:    true,
		},
		{
			description: "IsIPV6",
			target:      pointerOf("007e:89eb:9dc0:73cd:326f:503e:38d5:697e"),
			expected:    false,
		},
		{
			description: "IsNotIP",
			target:      pointerOf("string"),
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.IPv4Address(test.target)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: "{0} is not a valid IPv4 address",
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}

func TestIPv6Address(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *string
		expected    bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			expected:    false,
		},
		{
			description: "IsIPV4",
			target:      pointerOf("186.15.46.163"),
			expected:    false,
		},
		{
			description: "IsIPV6",
			target:      pointerOf("007e:89eb:9dc0:73cd:326f:503e:38d5:697e"),
			expected:    true,
		},
		{
			description: "IsNotIP",
			target:      pointerOf("string"),
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.IPv6Address(test.target)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: "{0} is not a valid IPv6 address",
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}

func TestDateTime(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *string
		layout      string
		expected    bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			layout:      time.RFC3339,
			expected:    false,
		},
		{
			description: "EmptyLayout",
			target:      nil,
			layout:      "",
			expected:    false,
		},
		{
			description: "IsRFC3339",
			target:      pointerOf("2023-10-20T15:04:05Z"),
			layout:      time.RFC3339,
			expected:    true,
		},
		{
			description: "IsRFC1123",
			target:      pointerOf("Fri, 20 Oct 2023 15:04:05 MST"),
			layout:      time.RFC1123,
			expected:    true,
		},
		{
			description: "IsNotDateTime",
			target:      pointerOf("string"),
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.DateTime(test.target, test.layout)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: fmt.Sprintf(`{0} is not a valid date time in required format (ex: "%s")`, test.layout),
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}
