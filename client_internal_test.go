package grpcsteps

import (
	"testing"

	"github.com/cucumber/godog"
	"github.com/nhatthm/grpcmock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/test/bufconn"

	"github.com/godogx/grpcsteps/internal/grpctest"
)

func TestClient_iRequestWithPayload_InvalidMethod(t *testing.T) {
	t.Parallel()

	err := NewClient().iRequestWithPayload("", nil)
	expected := `invalid grpc method`

	assert.EqualError(t, err, expected)
}

func TestClient_iRequestWithPayload_InvalidPayload(t *testing.T) {
	t.Parallel()

	err := NewClient(
		RegisterService(grpctest.RegisterItemServiceServer),
	).iRequestWithPayload("/grpctest.ItemService/GetItem", &godog.DocString{Content: "42"})
	expected := `json: cannot unmarshal number into Go value of type grpctest.GetItemRequest`

	assert.EqualError(t, err, expected)
}

func TestClient_iRequestWithTimeout_InvalidDuration(t *testing.T) {
	t.Parallel()

	err := NewClient().iRequestWithTimeout("not a timeout")
	expected := `time: invalid duration "not a timeout"`

	assert.EqualError(t, err, expected)
}

func TestClient_iShouldHaveResponseWithCode_InvalidCode(t *testing.T) {
	t.Parallel()

	err := NewClient().iShouldHaveResponseWithCode(`not a code`)
	expected := `invalid code: "\"NOT A CODE\""`

	assert.EqualError(t, err, expected)
}

func TestClient_iShouldHaveResponseWithCodeAndErrorMessage_InvalidCode(t *testing.T) {
	t.Parallel()

	err := NewClient().iShouldHaveResponseWithCodeAndErrorMessage(`not a code`, "")
	expected := `invalid code: "\"NOT A CODE\""`

	assert.EqualError(t, err, expected)
}

func TestClient_iShouldHaveResponseWithCodeAndErrorMessageFromDocString_InvalidCode(t *testing.T) {
	t.Parallel()

	err := NewClient().iShouldHaveResponseWithCodeAndErrorMessageFromDocString(`not a code`, &godog.DocString{})
	expected := `invalid code: "\"NOT A CODE\""`

	assert.EqualError(t, err, expected)
}

func TestToPayload(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		methodType     grpcmock.MethodType
		data           string
		expectedResult interface{}
		expectedError  string
	}{
		{
			scenario:      "invalid data for unary",
			methodType:    grpcmock.MethodTypeUnary,
			data:          "42",
			expectedError: "json: cannot unmarshal number into Go value of type grpctest.Item",
		},
		{
			scenario:      "invalid data for stream",
			methodType:    grpcmock.MethodTypeClientStream,
			data:          `{"id": 42}`,
			expectedError: "json: cannot unmarshal object into Go value of type []*grpctest.Item",
		},
		{
			scenario:       "unary payload",
			methodType:     grpcmock.MethodTypeUnary,
			data:           `{"id": 42}`,
			expectedResult: &grpctest.Item{Id: 42},
		},
		{
			scenario:       "server stream payload",
			methodType:     grpcmock.MethodTypeServerStream,
			data:           `{"id": 42}`,
			expectedResult: &grpctest.Item{Id: 42},
		},
		{
			scenario:       "client stream payload",
			methodType:     grpcmock.MethodTypeClientStream,
			data:           `[{"id": 42}]`,
			expectedResult: []*grpctest.Item{{Id: 42}},
		},
		{
			scenario:       "bidirectional stream payload",
			methodType:     grpcmock.MethodTypeBidirectionalStream,
			data:           `[{"id": 42}]`,
			expectedResult: []*grpctest.Item{{Id: 42}},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			result, err := toPayload(tc.methodType, &grpctest.Item{}, tc.data)

			assert.Equal(t, tc.expectedResult, result)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestToStatusCode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		code           string
		expectedResult codes.Code
		expectedError  string
	}{
		{
			scenario:       "invalid code",
			code:           "not a code",
			expectedResult: codes.Unknown,
			expectedError:  `invalid code: "\"NOT A CODE\""`,
		},
		{
			scenario:       "screaming snake case",
			code:           "DEADLINE_EXCEEDED",
			expectedResult: codes.DeadlineExceeded,
		},
		{
			scenario:       "camel case",
			code:           "DeadlineExceeded",
			expectedResult: codes.DeadlineExceeded,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			result, err := toStatusCode(tc.code)

			assert.Equal(t, tc.expectedResult, result)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestNewResponse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario   string
		methodtype grpcmock.MethodType
		expected   interface{}
	}{
		{
			scenario:   "unary",
			methodtype: grpcmock.MethodTypeUnary,
			expected:   &grpctest.Item{},
		},
		{
			scenario:   "client stream",
			methodtype: grpcmock.MethodTypeClientStream,
			expected:   &grpctest.Item{},
		},
		{
			scenario:   "server stream",
			methodtype: grpcmock.MethodTypeServerStream,
			expected:   &[]*grpctest.Item{},
		},
		{
			scenario:   "bidirectional",
			methodtype: grpcmock.MethodTypeBidirectionalStream,
			expected:   &[]*grpctest.Item{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual := newResponse(tc.methodtype, grpctest.Item{})

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestWithAddr(t *testing.T) {
	t.Parallel()

	addr := "127.0.0.1:9000"
	c := NewClient(RegisterService(grpctest.RegisterItemServiceServer, WithAddr(addr)))

	assert.Equal(t, addr, c.services["/grpctest.ItemService/GetItem"].Address)
}

func TestWithAddressProvider(t *testing.T) {
	t.Parallel()

	provider := bufconn.Listen(1024 * 1024)
	addr := "bufconn"

	defer provider.Close() // nolint: errcheck

	c := NewClient(RegisterService(grpctest.RegisterItemServiceServer, WithAddressProvider(provider)))

	assert.Equal(t, addr, c.services["/grpctest.ItemService/ListItems"].Address)
}
