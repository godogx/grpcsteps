package grpcsteps

import (
	"context"
	"net"

	"github.com/cucumber/godog"
	grpcReflect "github.com/nhatthm/grpcmock/reflect"
	"github.com/nhatthm/grpcmock/service"
	"google.golang.org/grpc"
)

// Service contains needed information to form a GRPC request.
type Service struct {
	service.Method

	Address     string
	DialOptions []grpc.DialOption
}

// ServiceOption sets up a service.
type ServiceOption func(s *Service)

// Client is a grpc client for godog.
type Client struct {
	services map[string]*Service

	defaultSvcOptions []ServiceOption
}

// ClientOption sets up a client.
type ClientOption func(s *Client)

func (c *Client) registerService(id string, svc interface{}, opts ...ServiceOption) {
	for _, method := range grpcReflect.FindServiceMethods(svc) {
		svc := &Service{
			Method: service.Method{
				ServiceName: id,
				MethodName:  method.Name,
				MethodType:  service.ToType(method.IsClientStream, method.IsServerStream),
				Input:       method.Input,
				Output:      method.Output,
			},
			Address: ":9090",
		}

		// Apply default options.
		for _, o := range c.defaultSvcOptions {
			o(svc)
		}

		// Apply inline options.
		for _, o := range opts {
			o(svc)
		}

		c.services[svc.FullName()] = svc
	}
}

// RegisterContext registers to godog scenario.
func (c *Client) RegisterContext(sc *godog.ScenarioContext) {
	sc.Step(`^I request(?: a)? (?:gRPC|GRPC|grpc)(?: method)? "([^"]*)" with payload:?$`, c.iRequestWithPayload)

	sc.Step(`^I should have(?: a)? (?:gRPC|GRPC|grpc) response with payload:?$`, c.iShouldHaveResponseWithResponse)
	sc.Step(`^I should have(?: a)? (?:gRPC|GRPC|grpc) response with code "([^"]*)"$`, c.iShouldHaveResponseWithCode)
	sc.Step(`^I should have(?: a)? (?:gRPC|GRPC|grpc) response with error (?:message )?"([^"]*)"$`, c.iShouldHaveResponseWithErrorMessage)
	sc.Step(`^I should have(?: a)? (?:gRPC|GRPC|grpc) response with code "([^"]*)" and error (?:message )?"([^"]*)"$`, c.iShouldHaveResponseWithCodeAndErrorMessage)
	sc.Step(`^I should have(?: a)? (?:gRPC|GRPC|grpc) response with error(?: message)?:$`, c.iShouldHaveResponseWithErrorMessageFromDocString)
	sc.Step(`^I should have(?: a)? (?:gRPC|GRPC|grpc) response with code "([^"]*)" and error(?: message)?:$`, c.iShouldHaveResponseWithCodeAndErrorMessageFromDocString)

	registerRequestPlanner(sc)
}

func (c *Client) iRequestWithPayload(ctx context.Context, method string, data *godog.DocString) (context.Context, error) {
	svc, ok := c.services[method]
	if !ok {
		return ctx, ErrInvalidGRPCMethod
	}

	payload, err := toPayload(svc.MethodType, svc.Input, data.Content)
	if err != nil {
		return nil, err
	}

	return newClientRequestPlannerContext(ctx, svc, payload), nil
}

func (c *Client) iShouldHaveResponseWithResponse(ctx context.Context, response *godog.DocString) error {
	return assertServerResponsePayload(clientRequestFromContext(ctx), response.Content)
}

func (c *Client) iShouldHaveResponseWithCode(ctx context.Context, codeValue string) error {
	code, err := toStatusCode(codeValue)
	if err != nil {
		return err
	}

	return assertServerResponseErrorCode(clientRequestFromContext(ctx), code)
}

func (c *Client) iShouldHaveResponseWithErrorMessage(ctx context.Context, err string) error {
	return assertServerResponseErrorMessage(clientRequestFromContext(ctx), err)
}

func (c *Client) iShouldHaveResponseWithCodeAndErrorMessage(ctx context.Context, codeValue, err string) error {
	if err := c.iShouldHaveResponseWithCode(ctx, codeValue); err != nil {
		return err
	}

	return c.iShouldHaveResponseWithErrorMessage(ctx, err)
}

func (c *Client) iShouldHaveResponseWithErrorMessageFromDocString(ctx context.Context, err *godog.DocString) error {
	return c.iShouldHaveResponseWithErrorMessage(ctx, err.Content)
}

func (c *Client) iShouldHaveResponseWithCodeAndErrorMessageFromDocString(ctx context.Context, codeValue string, err *godog.DocString) error {
	if err := c.iShouldHaveResponseWithCode(ctx, codeValue); err != nil {
		return err
	}

	return c.iShouldHaveResponseWithErrorMessage(ctx, err.Content)
}

// NewClient initiates a new grpc server extension for testing.
func NewClient(opts ...ClientOption) *Client {
	s := &Client{
		services: make(map[string]*Service),
	}

	for _, o := range opts {
		o(s)
	}

	return s
}

// RegisterServiceFromInstance registers a grpc server by its interface.
func RegisterServiceFromInstance(id string, svc interface{}, opts ...ServiceOption) ClientOption {
	return func(c *Client) {
		c.registerService(id, svc, opts...)
	}
}

// RegisterService registers a grpc server by its interface.
func RegisterService(registerFunc interface{}, opts ...ServiceOption) ClientOption {
	return func(c *Client) {
		serviceDesc, svc := grpcReflect.ParseRegisterFunc(registerFunc)

		c.registerService(serviceDesc.ServiceName, svc, opts...)
	}
}

// WithDefaultServiceOptions set default service options.
func WithDefaultServiceOptions(opts ...ServiceOption) ClientOption {
	return func(s *Client) {
		s.defaultSvcOptions = append(s.defaultSvcOptions, opts...)
	}
}

// AddrProvider provides a net address.
type AddrProvider interface {
	Addr() net.Addr
}

// WithAddressProvider sets service address.
func WithAddressProvider(p AddrProvider) ServiceOption {
	return WithAddr(p.Addr().String())
}

// WithAddr sets service address.
func WithAddr(addr string) ServiceOption {
	return func(s *Service) {
		s.Address = addr
	}
}

// WithDialOption adds a dial option.
func WithDialOption(o grpc.DialOption) ServiceOption {
	return func(s *Service) {
		s.DialOptions = append(s.DialOptions, o)
	}
}

// WithDialOptions sets dial options.
func WithDialOptions(opts ...grpc.DialOption) ServiceOption {
	return func(s *Service) {
		s.DialOptions = opts
	}
}
