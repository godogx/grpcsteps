package grpcsteps

import (
	"context"
	"fmt"
	"os"

	"github.com/cucumber/godog"
	"go.nhat.io/grpcmock"
	"go.nhat.io/grpcmock/planner"
	"go.nhat.io/grpcmock/service"
	"google.golang.org/grpc/codes"
)

// ExternalServiceManager is a grpc server for godog.
type ExternalServiceManager struct {
	servers map[string]*wrappedServer
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

func (m *ExternalServiceManager) receiveRequest(ctx context.Context, serviceID, method string, times uint, payload *string) (context.Context, error) {
	srv, found := m.servers[serviceID]
	if !found {
		//goland:noinspection GoErrorStringFormat
		return ctx, fmt.Errorf(
			"%w, did you forget to setup the grpc service %q?",
			ErrGRPCServiceNotFound, serviceID,
		)
	}

	r, err := srv.expect(method, times, payload)
	if err != nil {
		return ctx, err
	}

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
	return m.receiveRequest(ctx, service, method, uint(times), nil)
}

func (m *ExternalServiceManager) receiveRepeatedRequestsWithPayload(ctx context.Context, service string, times int, method, data string) (context.Context, error) {
	return m.receiveRequest(ctx, service, method, uint(times), &data)
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
	return m.receiveRequest(ctx, service, method, planner.UnlimitedTimes, nil)
}

func (m *ExternalServiceManager) receiveManyRequestsWithPayload(ctx context.Context, service, method, data string) (context.Context, error) {
	return m.receiveRequest(ctx, service, method, planner.UnlimitedTimes, &data)
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
	m.servers[id] = newServer(opts...)

	return m.servers[id].Address()
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
		servers: make(map[string]*wrappedServer),
	}
}

type wrappedServer struct {
	*grpcmock.Server
}

func (s *wrappedServer) expect(method string, times uint, payload *string) (expectation, error) {
	svc := grpcmock.FindServerMethod(s.Server, method)
	if svc == nil {
		return nil, fmt.Errorf("%w: %s", ErrGRPCMethodNotFound, method)
	}

	var expected expectation

	switch svc.MethodType {
	case service.TypeUnary:
		expected = &unaryExpectation{s.ExpectUnary(method)}

	case service.TypeClientStream:
		expected = &clientStreamExpectation{s.ExpectClientStream(method)}

	case service.TypeServerStream:
		expected = &serverStreamExpectation{s.ExpectServerStream(method)}

	case service.TypeBidirectionalStream:
		return nil, fmt.Errorf("%w: %s %s", ErrGRPCMethodNotSupported, svc.MethodType, method)
	}

	expected.Times(times)

	if payload != nil {
		expected.WithPayload(*payload)
	}

	return expected, nil
}

func newServer(opts ...grpcmock.ServerOption) *wrappedServer {
	return &wrappedServer{
		Server: grpcmock.NewServer(opts...),
	}
}

type expectation interface {
	WithPayload(in interface{})
	WithHeader(key string, value interface{})
	Return(v interface{})
	ReturnError(code codes.Code, msg string)
	Times(i uint)
}

type unaryExpectation struct {
	grpcmock.UnaryExpectation
}

func (e *unaryExpectation) WithPayload(in interface{}) {
	e.UnaryExpectation.WithPayload(in)
}

func (e *unaryExpectation) WithHeader(key string, value interface{}) {
	e.UnaryExpectation.WithHeader(key, value)
}

func (e *unaryExpectation) Return(v interface{}) {
	e.UnaryExpectation.Return(v)
}

func (e *unaryExpectation) ReturnError(code codes.Code, msg string) {
	e.UnaryExpectation.ReturnError(code, msg)
}

func (e *unaryExpectation) Times(i uint) {
	e.UnaryExpectation.Times(i)
}

type clientStreamExpectation struct {
	grpcmock.ClientStreamExpectation
}

func (e *clientStreamExpectation) WithPayload(in interface{}) {
	e.ClientStreamExpectation.WithPayload(in)
}

func (e *clientStreamExpectation) WithHeader(key string, value interface{}) {
	e.ClientStreamExpectation.WithHeader(key, value)
}

func (e *clientStreamExpectation) Return(v interface{}) {
	e.ClientStreamExpectation.Return(v)
}

func (e *clientStreamExpectation) ReturnError(code codes.Code, msg string) {
	e.ClientStreamExpectation.ReturnError(code, msg)
}

func (e *clientStreamExpectation) Times(i uint) {
	e.ClientStreamExpectation.Times(i)
}

type serverStreamExpectation struct {
	grpcmock.ServerStreamExpectation
}

func (e *serverStreamExpectation) WithPayload(in interface{}) {
	e.ServerStreamExpectation.WithPayload(in)
}

func (e *serverStreamExpectation) WithHeader(key string, value interface{}) {
	e.ServerStreamExpectation.WithHeader(key, value)
}

func (e *serverStreamExpectation) Return(v interface{}) {
	e.ServerStreamExpectation.Return(v)
}

func (e *serverStreamExpectation) ReturnError(code codes.Code, msg string) {
	e.ServerStreamExpectation.ReturnError(code, msg)
}

func (e *serverStreamExpectation) Times(i uint) {
	e.ServerStreamExpectation.Times(i)
}
