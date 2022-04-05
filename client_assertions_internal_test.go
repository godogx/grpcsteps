package grpcsteps

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAssertServerResponsePayload(t *testing.T) {
	t.Parallel()

	const expectedPayload = `{"name": "john"}`

	testCases := []struct {
		scenario      string
		expected      string
		request       clientRequestDoer
		expectedError string
	}{
		{
			scenario: "has error",
			request: func() ([]byte, error) {
				return nil, errors.New("request error")
			},
			expectedError: `an error occurred while send grpc request: request error`,
		},
		{
			scenario: "different payload",
			expected: expectedPayload,
			request: func() ([]byte, error) {
				return []byte(`{"name": "foobar"}`), nil
			},
			expectedError: `not equal:
 {
-  "name": "john"
+  "name": "foobar"
 }
`,
		},
		{
			scenario: "same payload",
			expected: expectedPayload,
			request: func() ([]byte, error) {
				return []byte(`{"name": "john"}`), nil
			},
		},
		{
			scenario: "ignore diff",
			expected: `{"name": "<ignore-diff>"}`,
			request: func() ([]byte, error) {
				return []byte(`{"name": "john"}`), nil
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			err := assertServerResponsePayload(tc.request, tc.expected)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestAssertServerResponseErrorCode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		expected      codes.Code
		request       clientRequestDoer
		expectedError string
	}{
		{
			scenario: "no error and expect OK",
			expected: codes.OK,
			request: func() ([]byte, error) {
				return nil, nil
			},
		},
		{
			scenario: "no error and expect failure",
			expected: codes.Internal,
			request: func() ([]byte, error) {
				return nil, nil
			},
			expectedError: `got no error, want "Internal"`,
		},
		{
			scenario: "got error and expect different",
			expected: codes.FailedPrecondition,
			request: func() ([]byte, error) {
				return nil, status.Error(codes.Internal, "internal server error")
			},
			expectedError: `got rpc error: code = Internal desc = internal server error, want "FailedPrecondition"`,
		},
		{
			scenario: "got expected error",
			expected: codes.Internal,
			request: func() ([]byte, error) {
				return nil, status.Error(codes.Internal, "internal server error")
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			err := assertServerResponseErrorCode(tc.request, tc.expected)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestAssertServerResponseErrorMessage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		expected      string
		request       clientRequestDoer
		expectedError string
	}{
		{
			scenario: "no error and expect OK",
			request: func() ([]byte, error) {
				return nil, nil
			},
		},
		{
			scenario: "no error and expect failure",
			expected: "Internal",
			request: func() ([]byte, error) {
				return nil, nil
			},
			expectedError: `got no error, want "Internal"`,
		},
		{
			scenario: "got error and expect different",
			expected: "server went away",
			request: func() ([]byte, error) {
				return nil, status.Error(codes.Internal, "internal server error")
			},
			expectedError: `unexpected error message, got "internal server error", want "server went away"`,
		},
		{
			scenario: "got expected error",
			expected: "internal server error",
			request: func() ([]byte, error) {
				return nil, status.Error(codes.Internal, "internal server error")
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			err := assertServerResponseErrorMessage(tc.request, tc.expected)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

type clientRequestDoer func() ([]byte, error)

func (d clientRequestDoer) Do() ([]byte, error) {
	return d()
}
