package grpcsteps

import (
	"context"
	"testing"

	"github.com/nhatthm/grpcmock/service"
	"github.com/stretchr/testify/assert"

	"github.com/godogx/grpcsteps/internal/grpctest"
)

func TestClientRequestInContext(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Case 1: no request in context.
	assert.Equal(t, missingClientRequest{}, clientRequestFromContext(ctx))

	// Case 2: request in context.
	ctx = clientRequestToContext(ctx, &clientRequestInvoker{})

	assert.Equal(t, &clientRequestInvoker{}, clientRequestFromContext(ctx))
}

func TestMissingClientRequest(t *testing.T) {
	t.Parallel()

	expected := "no client request in context, did you forget to setup a gprc request with `^I request(?: a)? (?:gRPC|GRPC|grpc)(?: method)? \"([^\"]*)\" with payload:?$` in the scenario?"

	r := missingClientRequest{}

	result, err := r.Do()

	assert.Nil(t, result)
	assert.EqualError(t, err, expected)
}

func TestNewServerOutput(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario   string
		methodtype service.Type
		expected   interface{}
	}{
		{
			scenario:   "unary",
			methodtype: service.TypeUnary,
			expected:   &grpctest.Item{},
		},
		{
			scenario:   "client stream",
			methodtype: service.TypeClientStream,
			expected:   &grpctest.Item{},
		},
		{
			scenario:   "server stream",
			methodtype: service.TypeServerStream,
			expected:   &[]*grpctest.Item{},
		},
		{
			scenario:   "bidirectional",
			methodtype: service.TypeBidirectionalStream,
			expected:   &[]*grpctest.Item{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual := newServerOutput(tc.methodtype, grpctest.Item{})

			assert.Equal(t, tc.expected, actual)
		})
	}
}
