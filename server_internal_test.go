package grpcsteps

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExternalServiceManager_ReceiveOneRequestWithPayloadFromFile_ReadFileError(t *testing.T) {
	t.Parallel()

	_, err := NewExternalServiceManager().
		receiveOneRequestWithPayloadFromFile(context.Background(), "item-service", "/grpctest.ItemService/GetItem", "not_found")

	expected := `open not_found: no such file or directory`

	assert.EqualError(t, err, expected)
}

func TestExternalServiceManager_ReceiveRepeatedRequestsWithPayloadFromFile_ReadFileError(t *testing.T) {
	t.Parallel()

	_, err := NewExternalServiceManager().
		receiveRepeatedRequestsWithPayloadFromFile(context.Background(), "item-service", 10, "/grpctest.ItemService/GetItem", "not_found")

	expected := `open not_found: no such file or directory`

	assert.EqualError(t, err, expected)
}

func TestExternalServiceManager_ReceiveManyRequestsWithPayloadFromFile_ReadFileError(t *testing.T) {
	t.Parallel()

	_, err := NewExternalServiceManager().
		receiveManyRequestsWithPayloadFromFile(context.Background(), "item-service", "/grpctest.ItemService/GetItem", "not_found")

	expected := `open not_found: no such file or directory`

	assert.EqualError(t, err, expected)
}

func TestExternalServiceManager_respondWithPayloadFromFile_ReadFileError(t *testing.T) {
	t.Parallel()

	err := NewExternalServiceManager().
		respondWithPayloadFromFile(context.Background(), "not_found")

	expected := `open not_found: no such file or directory`

	assert.EqualError(t, err, expected)
}

func TestExternalServiceManager_RespondWithError_InvalidCode(t *testing.T) {
	t.Parallel()

	err := NewExternalServiceManager().respondWithError(context.Background(), `not a code`, ``)
	expected := `invalid code: "\"NOT A CODE\""`

	assert.EqualError(t, err, expected)
}
