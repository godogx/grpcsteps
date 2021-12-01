package grpcsteps

import (
	"testing"

	"github.com/nhatthm/grpcmock/service"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"

	"github.com/godogx/grpcsteps/internal/grpctest"
)

func TestToPayload(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		methodType     service.Type
		data           string
		expectedResult interface{}
		expectedError  string
	}{
		{
			scenario:      "invalid data for unary",
			methodType:    service.TypeUnary,
			data:          "42",
			expectedError: "json: cannot unmarshal number into Go value of type grpctest.Item",
		},
		{
			scenario:      "invalid data for stream",
			methodType:    service.TypeClientStream,
			data:          `{"id": 42}`,
			expectedError: "json: cannot unmarshal object into Go value of type []*grpctest.Item",
		},
		{
			scenario:       "unary payload",
			methodType:     service.TypeUnary,
			data:           `{"id": 42}`,
			expectedResult: &grpctest.Item{Id: 42},
		},
		{
			scenario:       "server stream payload",
			methodType:     service.TypeServerStream,
			data:           `{"id": 42}`,
			expectedResult: &grpctest.Item{Id: 42},
		},
		{
			scenario:       "client stream payload",
			methodType:     service.TypeClientStream,
			data:           `[{"id": 42}]`,
			expectedResult: []*grpctest.Item{{Id: 42}},
		},
		{
			scenario:       "bidirectional stream payload",
			methodType:     service.TypeBidirectionalStream,
			data:           `[{"id": 42}]`,
			expectedResult: []*grpctest.Item{{Id: 42}},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			result, err := toPayload(tc.methodType, &grpctest.Item{}, tc.data)

			assert.Equal(t, tc.expectedResult, result)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestToStatusCode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		code           string
		expectedResult codes.Code
		expectedError  string
	}{
		{
			scenario:       "invalid code",
			code:           "not a code",
			expectedResult: codes.Unknown,
			expectedError:  `invalid code: "\"NOT A CODE\""`,
		},
		{
			scenario:       "screaming snake case",
			code:           "DEADLINE_EXCEEDED",
			expectedResult: codes.DeadlineExceeded,
		},
		{
			scenario:       "camel case",
			code:           "DeadlineExceeded",
			expectedResult: codes.DeadlineExceeded,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			result, err := toStatusCode(tc.code)

			assert.Equal(t, tc.expectedResult, result)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestToUpperSnakeCase(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		value    string
		expected string
	}{
		{
			scenario: "PascalCase",
			value:    "NotImplemented",
			expected: "NOT_IMPLEMENTED",
		},
		{
			scenario: "camelCase",
			value:    "notImplemented",
			expected: "NOT_IMPLEMENTED",
		},
		{
			scenario: "snake_case",
			value:    "not_implemented",
			expected: "NOT_IMPLEMENTED",
		},
		{
			scenario: "UPPER_SNAKE_CASE",
			value:    "NOT_IMPLEMENTED",
			expected: "NOT_IMPLEMENTED",
		},
		{
			scenario: "kebab-case is not supported",
			value:    "not-implemented",
			expected: "NOT-IMPLEMENTED",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, toUpperSnakeCase(tc.value))
		})
	}
}
