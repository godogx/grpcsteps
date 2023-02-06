package grpcsteps

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerRequestReflectorPlanner_WithTimeout(t *testing.T) {
	t.Parallel()

	p := newServerRequestPlanner((*unaryExpectation)(nil))

	err := p.WithTimeout(0)

	expected := `grpc service request does not have timeout`

	assert.EqualError(t, err, expected)
}

func TestServerRequestPlannerInContext(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Case 1: no planner in context.
	assert.Equal(t, missingServerRequestPlanner{}, serverRequestPlannerFromContext(ctx))

	// Case 2: request in context.
	ctx = newServerRequestPlannerContext(ctx, nil)

	assert.Equal(t, &serverRequestReflectorPlanner{}, serverRequestPlannerFromContext(ctx))
}

func TestMissingServerRequestPlanner(t *testing.T) {
	t.Parallel()

	expected := `no service request in context, did you forget to setup a gprc request in the scenario?

For example:

        When "item-service" receives a grpc request "/grpctest.ItemService/GetItem" with payload:
        """
        {
            "id": 42
        }
        """
`

	p := missingServerRequestPlanner{}

	assert.EqualError(t, p.WithHeader("", nil), expected)
	assert.EqualError(t, p.WithTimeout(0), expected)
	assert.EqualError(t, p.Return(""), expected)
	assert.EqualError(t, p.ReturnError(0, ""), expected)
}
