package grpcsteps

import (
	"context"
	"fmt"
	"time"

	"github.com/cucumber/godog"
)

// ErrNoRequestPlannerInContext indicates that there is no request planner in context.
const ErrNoRequestPlannerInContext err = "no request planner in context"

type requestPlanner interface {
	WithHeader(header string, value interface{}) error
	WithTimeout(d time.Duration) error
}

func registerRequestPlanner(sc *godog.ScenarioContext) {
	sc.Step(`^The (?:gRPC|GRPC|grpc) request has(?: a)? header "([^"]*): ([^"]*)"$`, planRequestWithHeader)
	sc.Step(`^The (?:gRPC|GRPC|grpc) request timeout is "([^"]*)"$`, planRequestWithTimeout)
}

func planRequestWithHeader(ctx context.Context, header, value string) error {
	return requestPlannerFromContext(ctx).WithHeader(header, value)
}

func planRequestWithTimeout(ctx context.Context, t string) error {
	timeout, err := time.ParseDuration(t)
	if err != nil {
		return err
	}

	return requestPlannerFromContext(ctx).WithTimeout(timeout)
}

type requestCtxKey struct{}

type requestPlannerCtxKey struct{}

func requestPlannerFromContext(ctx context.Context) requestPlanner {
	p, ok := ctx.Value(requestPlannerCtxKey{}).(requestPlanner)
	if !ok {
		return missingRequestPlanner{}
	}

	return p
}

func requestPlannerToContext(ctx context.Context, r requestPlanner) context.Context {
	return context.WithValue(ctx, requestPlannerCtxKey{}, r)
}

type missingRequestPlanner struct{}

func (missingRequestPlanner) WithHeader(string, interface{}) error {
	return missingRequestPlannerErr()
}

func (missingRequestPlanner) WithTimeout(time.Duration) error {
	return missingRequestPlannerErr()
}

func missingRequestPlannerErr() error {
	//goland:noinspection GoErrorStringFormat
	return fmt.Errorf(
		"%w, did you forget to setup a gprc request in the scenario?",
		ErrNoRequestPlannerInContext,
	)
}
