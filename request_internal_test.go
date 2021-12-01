package grpcsteps

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlanRequestWithTimeout(t *testing.T) {
	t.Parallel()

	err := planRequestWithTimeout(context.Background(), "not a timeout")
	expected := `time: invalid duration "not a timeout"`

	assert.EqualError(t, err, expected)
}

func TestRequestPlannerInContext(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Case 1: no planner in context.
	assert.Equal(t, missingRequestPlanner{}, requestPlannerFromContext(ctx))

	// Case 2: planner in context.
	ctx = requestPlannerToContext(ctx, clientRequestPlanner{})

	assert.Equal(t, clientRequestPlanner{}, requestPlannerFromContext(ctx))
}

func TestMissingRequestPlanner(t *testing.T) {
	t.Parallel()

	expected := `no request planner in context, did you forget to setup a gprc request in the scenario?`

	p := missingRequestPlanner{}

	assert.EqualError(t, p.WithHeader("", nil), expected)
	assert.EqualError(t, p.WithTimeout(0), expected)
}
