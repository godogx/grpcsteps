package grpcsteps

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
)

// ErrNoServiceRequestInContext indicates that there is no service request in context.
const ErrNoServiceRequestInContext err = "no service request in context"

type serverRequestPlanner interface {
	requestPlanner

	Return(payload string) error
	ReturnError(code codes.Code, message string) error
}

type serverRequestReflectorPlanner struct {
	expected expectation
}

func (s *serverRequestReflectorPlanner) WithHeader(header string, value interface{}) error {
	s.expected.WithHeader(header, value)

	return nil
}

func (s *serverRequestReflectorPlanner) WithTimeout(time.Duration) error {
	return fmt.Errorf("grpc service request does not have timeout") // nolint: goerr113
}

func (s *serverRequestReflectorPlanner) Return(payload string) error { // nolint: unparam
	s.expected.Return(payload)

	return nil
}

func (s *serverRequestReflectorPlanner) ReturnError(code codes.Code, message string) error { // nolint: unparam
	s.expected.ReturnError(code, message)

	return nil
}

func newServerRequestPlanner(expected expectation) *serverRequestReflectorPlanner {
	return &serverRequestReflectorPlanner{
		expected: expected,
	}
}

func serverRequestPlannerFromContext(ctx context.Context) serverRequestPlanner {
	r, ok := ctx.Value(requestPlannerCtxKey{}).(serverRequestPlanner)
	if !ok {
		return missingServerRequestPlanner{}
	}

	return r
}

func newServerRequestPlannerContext(ctx context.Context, expected expectation) context.Context {
	return requestPlannerToContext(ctx, newServerRequestPlanner(expected))
}

type missingServerRequestPlanner struct{}

func (missingServerRequestPlanner) WithHeader(string, interface{}) error {
	return missingServerRequestPlannerErr()
}

func (missingServerRequestPlanner) WithTimeout(time.Duration) error {
	return missingServerRequestPlannerErr()
}

func (missingServerRequestPlanner) Return(string) error {
	return missingServerRequestPlannerErr()
}

func (missingServerRequestPlanner) ReturnError(codes.Code, string) error {
	return missingServerRequestPlannerErr()
}

func missingServerRequestPlannerErr() error {
	//goland:noinspection GoErrorStringFormat
	return fmt.Errorf(
		"%w, did you forget to setup a gprc request in the scenario?\n\nFor example:\n%s",
		ErrNoServiceRequestInContext,
		`
        When "item-service" receives a grpc request "/grpctest.ItemService/GetItem" with payload:
        """
        {
            "id": 42
        }
        """
`,
	)
}
