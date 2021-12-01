package grpcsteps

import (
	"context"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/test/bufconn"

	"github.com/godogx/grpcsteps/internal/grpctest"
)

func TestClient_iRequestWithPayload_InvalidMethod(t *testing.T) {
	t.Parallel()

	_, err := NewClient().iRequestWithPayload(context.Background(), "", nil)
	expected := `invalid grpc method`

	assert.EqualError(t, err, expected)
}

func TestClient_iRequestWithPayload_InvalidPayload(t *testing.T) {
	t.Parallel()

	_, err := NewClient(
		RegisterService(grpctest.RegisterItemServiceServer),
	).iRequestWithPayload(context.Background(), "/grpctest.ItemService/GetItem", &godog.DocString{Content: "42"})
	expected := `json: cannot unmarshal number into Go value of type grpctest.GetItemRequest`

	assert.EqualError(t, err, expected)
}

func TestClient_iShouldHaveResponseWithCode_InvalidCode(t *testing.T) {
	t.Parallel()

	err := NewClient().iShouldHaveResponseWithCode(context.Background(), `not a code`)
	expected := `invalid code: "\"NOT A CODE\""`

	assert.EqualError(t, err, expected)
}

func TestClient_iShouldHaveResponseWithCodeAndErrorMessage_InvalidCode(t *testing.T) {
	t.Parallel()

	err := NewClient().iShouldHaveResponseWithCodeAndErrorMessage(context.Background(), `not a code`, "")
	expected := `invalid code: "\"NOT A CODE\""`

	assert.EqualError(t, err, expected)
}

func TestClient_iShouldHaveResponseWithCodeAndErrorMessageFromDocString_InvalidCode(t *testing.T) {
	t.Parallel()

	err := NewClient().iShouldHaveResponseWithCodeAndErrorMessageFromDocString(context.Background(), `not a code`, &godog.DocString{})
	expected := `invalid code: "\"NOT A CODE\""`

	assert.EqualError(t, err, expected)
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
