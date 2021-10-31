package grpcsteps

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/cucumber/godog"
	"github.com/nhatthm/grpcmock"
	"github.com/nhatthm/grpcmock/invoker"
	grpcReflect "github.com/nhatthm/grpcmock/reflect"
	"github.com/stretchr/testify/assert"
	"github.com/swaggest/assertjson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

type clientExpect struct {
	invoker     *invoker.Invoker
	response    interface{}
	responseErr error
	isPopulated bool
}

func (e *clientExpect) doRequest() (interface{}, error) {
	if e.isPopulated {
		return e.response, e.responseErr
	}

	e.isPopulated = true
	e.responseErr = e.invoker.Invoke(context.Background())

	return e.response, e.responseErr
}

func (e *clientExpect) assertResponse(expected string) error {
	actual, err := e.doRequest()
	if err != nil {
		return err
	}

	t := &testT{}

	assertjson.EqualMarshal(t, []byte(expected), actual)

	return t.error
}

func (e *clientExpect) assertErrorCode(expected codes.Code) error {
	_, err := e.doRequest()
	if err == nil {
		if expected != codes.OK {
			return fmt.Errorf("got no error, want %q", expected.String()) // nolint: goerr113
		}

		return nil
	}

	actual := statusError(err)
	t := &testT{}

	assert.Equalf(t, expected, actual.Code(), "unexpected error code, got %q, want %q", actual.Code(), expected)

	return t.error
}

func (e *clientExpect) assertErrorMessage(expected string) error {
	_, err := e.doRequest()
	if err == nil {
		if expected != "" {
			return fmt.Errorf("got no error, want %q", expected) // nolint: goerr113
		}

		return nil
	}

	actual := statusError(err)
	t := &testT{}

	assert.Equalf(t, actual.Message(), expected, "unexpected error message, got %q, want %q", actual.Message(), expected)

	return t.error
}

// Service contains needed information to form a GRPC request.
type Service struct {
	grpcmock.ServiceMethod

	Address     string
	DialOptions []grpc.DialOption
}

// ServiceOption sets up a service.
type ServiceOption func(s *Service)

// Client is a grpc client for godog.
type Client struct {
	services map[string]*Service
	expect   *clientExpect

	defaultSvcOptions []ServiceOption
}

// ClientOption sets up a client.
type ClientOption func(s *Client)

func (c *Client) registerService(id string, svc interface{}, opts ...ServiceOption) {
	for _, method := range grpcReflect.FindServiceMethods(svc) {
		svc := &Service{
			ServiceMethod: grpcmock.ServiceMethod{
				ServiceName: id,
				MethodName:  method.Name,
				MethodType:  grpcmock.ToMethodType(method.IsClientStream, method.IsServerStream),
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
func (c *Client) RegisterContext(ctx *godog.ScenarioContext) {
	ctx.Before(func(context.Context, *godog.Scenario) (context.Context, error) {
		c.expect = nil

		return nil, nil
	})

	ctx.Step(`^I request(?: a)? (?:GRPC|grpc)(?: method)? "([^"]*)" with payload:?$`, c.iRequestWithPayload)
	ctx.Step(`^The (?:GRPC|grpc) request has(?: a)? header "([^"]*): ([^"]*)"$`, c.iRequestWithHeader)
	ctx.Step(`^The (?:GRPC|grpc) request timeout is "([^"]*)"$`, c.iRequestWithTimeout)

	ctx.Step(`^I should have(?: a)? (?:GRPC|grpc) response with payload:?$`, c.iShouldHaveResponseWithResponse)
	ctx.Step(`^I should have(?: a)? (?:GRPC|grpc) response with code "([^"]*)"$`, c.iShouldHaveResponseWithCode)
	ctx.Step(`^I should have(?: a)? (?:GRPC|grpc) response with error (?:message )?"([^"]*)"$`, c.iShouldHaveResponseWithErrorMessage)
	ctx.Step(`^I should have(?: a)? (?:GRPC|grpc) response with code "([^"]*)" and error (?:message )?"([^"]*)"$`, c.iShouldHaveResponseWithCodeAndErrorMessage)
}

func (c *Client) iRequestWithPayload(method string, params *godog.DocString) error {
	svc, ok := c.services[method]
	if !ok {
		return ErrInvalidGRPCMethod
	}

	e, err := newExpect(svc, params.Content)
	if err != nil {
		return err
	}

	c.expect = e

	return nil
}

func (c *Client) iRequestWithHeader(header, value string) error {
	c.expect.invoker.WithInvokeOption(grpcmock.WithHeader(header, value))

	return nil
}

func (c *Client) iRequestWithTimeout(t string) error {
	timeout, err := time.ParseDuration(t)
	if err != nil {
		return err
	}

	c.expect.invoker.WithTimeout(timeout)

	return nil
}

func (c *Client) iShouldHaveResponseWithResponse(response *godog.DocString) error {
	return c.expect.assertResponse(response.Content)
}

func (c *Client) iShouldHaveResponseWithCode(codeValue string) error {
	code, err := toStatusCode(codeValue)
	if err != nil {
		return err
	}

	return c.expect.assertErrorCode(code)
}

func (c *Client) iShouldHaveResponseWithErrorMessage(err string) error {
	return c.expect.assertErrorMessage(err)
}

func (c *Client) iShouldHaveResponseWithCodeAndErrorMessage(codeValue, err string) error {
	if err := c.iShouldHaveResponseWithCode(codeValue); err != nil {
		return err
	}

	return c.expect.assertErrorMessage(err)
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

func newExpect(svc *Service, data string) (*clientExpect, error) {
	payload, err := toPayload(svc.MethodType, svc.Input, data)
	if err != nil {
		return nil, err
	}

	response := newResponse(svc.MethodType, svc.Output)
	opts := []invoker.Option{
		invoker.WithAddress(svc.Address),
	}

	switch svc.MethodType {
	case grpcmock.MethodTypeBidirectionalStream:
		opts = append(opts, invoker.WithBidirectionalStreamHandler(grpcmock.SendAndRecvAll(payload, response)))

	case grpcmock.MethodTypeClientStream:
		opts = append(opts, invoker.WithInputStreamHandler(grpcmock.SendAll(payload)),
			invoker.WithOutput(response),
		)

	case grpcmock.MethodTypeServerStream:
		opts = append(opts, invoker.WithInput(payload),
			invoker.WithOutputStreamHandler(grpcmock.RecvAll(response)),
		)

	case grpcmock.MethodTypeUnary:
		fallthrough
	default:
		opts = append(opts, invoker.WithInput(payload),
			invoker.WithOutput(response),
		)
	}

	opts = append(opts, invoker.WithDialOptions(svc.DialOptions...))

	i := invoker.New(svc.ServiceMethod, opts...)

	return &clientExpect{
		invoker:     i,
		response:    response,
		responseErr: nil,
	}, nil
}

func newResponse(methodType grpcmock.MethodType, out interface{}) interface{} {
	result := reflect.New(grpcReflect.UnwrapType(out))

	if grpcmock.IsMethodServerStream(methodType) ||
		grpcmock.IsMethodBidirectionalStream(methodType) {
		value := reflect.MakeSlice(reflect.SliceOf(result.Type()), 0, 0)
		result = reflect.New(value.Type())

		result.Elem().Set(value)
	}

	return result.Interface()
}

func toPayload(methodType grpcmock.MethodType, in interface{}, data string) (interface{}, error) {
	result := reflect.New(grpcReflect.UnwrapType(in))
	isSlice := grpcmock.IsMethodClientStream(methodType) || grpcmock.IsMethodBidirectionalStream(methodType)

	if isSlice {
		result = reflect.MakeSlice(reflect.SliceOf(result.Type()), 0, 0)
		result = reflect.New(result.Type())
	}

	if err := json.Unmarshal([]byte(data), result.Interface()); err != nil {
		return nil, err
	}

	if isSlice {
		return result.Elem().Interface(), nil
	}

	return result.Interface(), nil
}

func toStatusCode(data string) (codes.Code, error) {
	data = fmt.Sprintf("%q", toUpperSnakeCase(data))

	var code codes.Code

	if err := code.UnmarshalJSON([]byte(data)); err != nil {
		return codes.Unknown, err
	}

	return code, nil
}

func toUpperSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")

	return strings.ToUpper(snake)
}

func statusError(err error) *status.Status {
	return status.Convert(err)
}
