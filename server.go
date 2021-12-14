package grpcsteps

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/cucumber/godog"
	"github.com/nhatthm/grpcmock"
	"github.com/nhatthm/grpcmock/request"
	"github.com/nhatthm/grpcmock/service"
	"google.golang.org/grpc/codes"
)

const (
	methodServerExpect       = "Expect"
	methodRequestWithPayload = "WithPayload"
	methodRequestWithHeader  = "WithHeader"
	methodRequestReturn      = "Return"
	methodRequestReturnError = "ReturnError"
	methodRequestTimes       = "Times"
)

// ExternalServiceManager is a grpc server for godog.
type ExternalServiceManager struct {
	servers map[string]*grpcmock.Server
}

// RegisterContext registers to godog scenario.
func (m *ExternalServiceManager) RegisterContext(sc *godog.ScenarioContext) {
	sc.Before(func(context.Context, *godog.Scenario) (context.Context, error) {
		m.resetExpectations()

		return nil, nil
	})

	sc.After(func(ctx context.Context, _ *godog.Scenario, err error) (context.Context, error) {
		if err != nil {
			return ctx, nil // nolint: nilerr
		}

		return ctx, m.assertExpectationsWereMet()
	})

	sc.Step(`^"([^"]*)" receives [a1] (?:gRPC|GRPC|grpc) request "([^"]*)"$`, m.receiveOneRequestWithoutPayload)
	sc.Step(`^"([^"]*)" receives [a1] (?:gRPC|GRPC|grpc) request "([^"]*)" with payload:$`, m.receiveOneRequestWithPayloadFromDocString)
	sc.Step(`^"([^"]*)" receives [a1] (?:gRPC|GRPC|grpc) request "([^"]*)" with payload from file "([^"]+)"$`, m.receiveOneRequestWithPayloadFromFile)
	sc.Step(`^"([^"]*)" receives [a1] (?:gRPC|GRPC|grpc) request "([^"]*)" with payload from file:$`, m.receiveOneRequestWithPayloadFromFileDocString)

	sc.Step(`^"([^"]*)" receives ([0-9]+) (?:gRPC|GRPC|grpc) requests "([^"]*)"$`, m.receiveRepeatedRequestsWithoutPayload)
	sc.Step(`^"([^"]*)" receives ([0-9]+) (?:gRPC|GRPC|grpc) requests "([^"]*)" with payload:$`, m.receiveRepeatedRequestsWithPayloadFromDocString)
	sc.Step(`^"([^"]*)" receives ([0-9]+) (?:gRPC|GRPC|grpc) requests "([^"]*)" with payload from file "([^"]+)"$`, m.receiveRepeatedRequestsWithPayloadFromFile)
	sc.Step(`^"([^"]*)" receives ([0-9]+) (?:gRPC|GRPC|grpc) requests "([^"]*)" with payload from file:$`, m.receiveRepeatedRequestsWithPayloadFromFileDocString)

	sc.Step(`^"([^"]*)" receives (?:some|many|several) (?:gRPC|GRPC|grpc) requests "([^"]*)"$`, m.receiveManyRequestsWithoutPayload)
	sc.Step(`^"([^"]*)" receives (?:some|many|several) (?:gRPC|GRPC|grpc) requests "([^"]*)" with payload:$`, m.receiveManyRequestsWithPayloadFromDocString)
	sc.Step(`^"([^"]*)" receives (?:some|many|several) (?:gRPC|GRPC|grpc) requests "([^"]*)" with payload from file "([^"]+)"$`, m.receiveManyRequestsWithPayloadFromFile)
	sc.Step(`^"([^"]*)" receives (?:some|many|several) (?:gRPC|GRPC|grpc) requests "([^"]*)" with payload from file:$`, m.receiveManyRequestsWithPayloadFromFileDocString)

	sc.Step(`^[tT]he (?:gRPC|GRPC|grpc) service responds with payload:?$`, m.respondWithPayloadFromDocString)
	sc.Step(`^[tT]he (?:gRPC|GRPC|grpc) service responds with payload from file "([^"]+)"$`, m.respondWithPayloadFromFile)
	sc.Step(`^[tT]he (?:gRPC|GRPC|grpc) service responds with payload from file:$`, m.respondWithPayloadFromFileDocString)
	sc.Step(`^[tT]he (?:gRPC|GRPC|grpc) service responds with code "([^"]*)"$`, m.respondWithErrorCode)
	sc.Step(`^[tT]he (?:gRPC|GRPC|grpc) service responds with error (?:message )?"([^"]*)"$`, m.respondWithErrorMessage)
	sc.Step(`^[tT]he (?:gRPC|GRPC|grpc) service responds with error(?: message)?:$`, m.respondWithErrorMessageFromDocString)
	sc.Step(`^[tT]he (?:gRPC|GRPC|grpc) service responds with code "([^"]*)" and error (?:message )?"([^"]*)"$`, m.respondWithError)
	sc.Step(`^[tT]he (?:gRPC|GRPC|grpc) service responds with code "([^"]*)" and error(?: message)?:$`, m.respondWithErrorFromDocString)

	registerRequestPlanner(sc)
}

func (m *ExternalServiceManager) receiveRequest(ctx context.Context, serviceID, method string, times request.RepeatedTime, payload *string) (context.Context, error) {
	srv, found := m.servers[serviceID]
	if !found {
		//goland:noinspection GoErrorStringFormat
		return ctx, fmt.Errorf(
			"%w, did you forget to setup the grpc service %q?",
			ErrGRPCServiceNotFound, serviceID,
		)
	}

	svc := grpcmock.FindServerMethod(srv, method)
	if svc == nil {
		return ctx, fmt.Errorf("%w: %s", ErrGRPCMethodNotFound, method)
	}

	if service.IsMethodBidirectionalStream(svc.MethodType) {
		return ctx, fmt.Errorf("%w: %s %s", ErrGRPCMethodNotSupported, svc.MethodType, method)
	}

	r := expectServerRequest(srv, svc, times, payload)

	return newServerRequestPlannerContext(ctx, r), nil
}

func (m *ExternalServiceManager) receiveOneRequestWithoutPayload(ctx context.Context, service, method string) (context.Context, error) {
	return m.receiveRequest(ctx, service, method, 1, nil)
}

func (m *ExternalServiceManager) receiveOneRequestWithPayload(ctx context.Context, service, method, data string) (context.Context, error) {
	return m.receiveRequest(ctx, service, method, 1, &data)
}

func (m *ExternalServiceManager) receiveOneRequestWithPayloadFromDocString(ctx context.Context, service, method string, payload *godog.DocString) (context.Context, error) {
	return m.receiveOneRequestWithPayload(ctx, service, method, payload.Content)
}

func (m *ExternalServiceManager) receiveOneRequestWithPayloadFromFile(ctx context.Context, service, method, path string) (context.Context, error) {
	payload, err := os.ReadFile(path) // nolint: gosec
	if err != nil {
		return ctx, err
	}

	return m.receiveOneRequestWithPayload(ctx, service, method, string(payload))
}

func (m *ExternalServiceManager) receiveOneRequestWithPayloadFromFileDocString(ctx context.Context, service, method string, path *godog.DocString) (context.Context, error) {
	return m.receiveOneRequestWithPayloadFromFile(ctx, service, method, path.Content)
}

func (m *ExternalServiceManager) receiveRepeatedRequestsWithoutPayload(ctx context.Context, service string, times int, method string) (context.Context, error) {
	return m.receiveRequest(ctx, service, method, request.RepeatedTime(times), nil)
}

func (m *ExternalServiceManager) receiveRepeatedRequestsWithPayload(ctx context.Context, service string, times int, method, data string) (context.Context, error) {
	return m.receiveRequest(ctx, service, method, request.RepeatedTime(times), &data)
}

func (m *ExternalServiceManager) receiveRepeatedRequestsWithPayloadFromDocString(ctx context.Context, service string, times int, method string, payload *godog.DocString) (context.Context, error) {
	return m.receiveRepeatedRequestsWithPayload(ctx, service, times, method, payload.Content)
}

func (m *ExternalServiceManager) receiveRepeatedRequestsWithPayloadFromFile(ctx context.Context, service string, times int, method, path string) (context.Context, error) {
	payload, err := os.ReadFile(path) // nolint: gosec
	if err != nil {
		return ctx, err
	}

	return m.receiveRepeatedRequestsWithPayload(ctx, service, times, method, string(payload))
}

func (m *ExternalServiceManager) receiveRepeatedRequestsWithPayloadFromFileDocString(ctx context.Context, service string, times int, method string, path *godog.DocString) (context.Context, error) {
	return m.receiveRepeatedRequestsWithPayloadFromFile(ctx, service, times, method, path.Content)
}

func (m *ExternalServiceManager) receiveManyRequestsWithoutPayload(ctx context.Context, service, method string) (context.Context, error) {
	return m.receiveRequest(ctx, service, method, request.UnlimitedTimes, nil)
}

func (m *ExternalServiceManager) receiveManyRequestsWithPayload(ctx context.Context, service, method, data string) (context.Context, error) {
	return m.receiveRequest(ctx, service, method, request.UnlimitedTimes, &data)
}

func (m *ExternalServiceManager) receiveManyRequestsWithPayloadFromDocString(ctx context.Context, service, method string, payload *godog.DocString) (context.Context, error) {
	return m.receiveManyRequestsWithPayload(ctx, service, method, payload.Content)
}

func (m *ExternalServiceManager) receiveManyRequestsWithPayloadFromFile(ctx context.Context, service, method, path string) (context.Context, error) {
	payload, err := os.ReadFile(path) // nolint: gosec
	if err != nil {
		return ctx, err
	}

	return m.receiveManyRequestsWithPayload(ctx, service, method, string(payload))
}

func (m *ExternalServiceManager) receiveManyRequestsWithPayloadFromFileDocString(ctx context.Context, service, method string, path *godog.DocString) (context.Context, error) {
	return m.receiveManyRequestsWithPayloadFromFile(ctx, service, method, path.Content)
}

func (m *ExternalServiceManager) respondWithPayload(ctx context.Context, payload string) error {
	return serverRequestPlannerFromContext(ctx).Return(payload)
}

func (m *ExternalServiceManager) respondWithPayloadFromDocString(ctx context.Context, payload *godog.DocString) error {
	return m.respondWithPayload(ctx, payload.Content)
}

func (m *ExternalServiceManager) respondWithPayloadFromFile(ctx context.Context, path string) error {
	payload, err := os.ReadFile(path) // nolint: gosec
	if err != nil {
		return err
	}

	return m.respondWithPayload(ctx, string(payload))
}

func (m *ExternalServiceManager) respondWithPayloadFromFileDocString(ctx context.Context, path *godog.DocString) error {
	return m.respondWithPayloadFromFile(ctx, path.Content)
}

func (m *ExternalServiceManager) respondWithError(ctx context.Context, codeValue string, message string) error {
	code, err := toStatusCode(codeValue)
	if err != nil {
		return err
	}

	return serverRequestPlannerFromContext(ctx).ReturnError(code, message)
}

func (m *ExternalServiceManager) respondWithErrorFromDocString(ctx context.Context, codeValue string, message *godog.DocString) error {
	return m.respondWithError(ctx, codeValue, message.Content)
}

func (m *ExternalServiceManager) respondWithErrorCode(ctx context.Context, codeValue string) error {
	return m.respondWithError(ctx, codeValue, "")
}

func (m *ExternalServiceManager) respondWithErrorMessage(ctx context.Context, message string) error {
	return m.respondWithError(ctx, codes.Internal.String(), message)
}

func (m *ExternalServiceManager) respondWithErrorMessageFromDocString(ctx context.Context, message *godog.DocString) error {
	return m.respondWithErrorMessage(ctx, message.Content)
}

func (m *ExternalServiceManager) resetExpectations() {
	for _, srv := range m.servers {
		srv.ResetExpectations()
	}
}

// AddService starts a new service and returns the server address for client to connect.
func (m *ExternalServiceManager) AddService(id string, opts ...grpcmock.ServerOption) string {
	srv := grpcmock.NewServer(opts...)

	m.servers[id] = srv

	return srv.Address()
}

// Close closes the server.
func (m *ExternalServiceManager) Close() {
	for _, srv := range m.servers {
		_ = srv.Close() // nolint: errcheck
	}
}

func (m *ExternalServiceManager) assertExpectationsWereMet() error {
	for _, srv := range m.servers {
		if err := srv.ExpectationsWereMet(); err != nil {
			return err
		}
	}

	return nil
}

// NewExternalServiceManager initiates a new external service manager for testing.
func NewExternalServiceManager() *ExternalServiceManager {
	return &ExternalServiceManager{
		servers: make(map[string]*grpcmock.Server),
	}
}

func callMethod(obj interface{}, method string, args ...interface{}) []reflect.Value {
	callArgs := make([]reflect.Value, len(args))

	for i, arg := range args {
		callArgs[i] = reflect.ValueOf(arg)
	}

	return reflect.ValueOf(obj).
		MethodByName(method).
		Call(callArgs)
}

func expectServerRequest(srv *grpcmock.Server, svc *service.Method, times request.RepeatedTime, payload *string) request.Request {
	method := fmt.Sprintf("%s%s", methodServerExpect, svc.MethodType)

	result := callMethod(srv, method, svc.FullName())
	r := result[0].Interface().(request.Request) // nolint: errcheck

	if payload != nil {
		expectServerRequestWithPayload(r, *payload)
	}

	callMethod(r, methodRequestTimes, times)

	return r
}

func expectServerRequestWithPayload(r request.Request, payload string) {
	callMethod(r, methodRequestWithPayload, payload)
}

func expectServerRequestWithHeader(r request.Request, header string, value interface{}) {
	callMethod(r, methodRequestWithHeader, header, value)
}

func setServerRequestReturn(r request.Request, payload string) {
	callMethod(r, methodRequestReturn, payload)
}

func setServerRequestReturnError(r request.Request, code codes.Code, message string) {
	callMethod(r, methodRequestReturnError, code, message)
}
