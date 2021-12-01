package grpcsteps

import (
	"fmt"

	"github.com/swaggest/assertjson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func assertServerResponsePayload(req clientRequest, expected string) error {
	actual, err := req.Do()
	if err != nil {
		return fmt.Errorf("an error occurred while send grpc request: %w", err)
	}

	return assertjson.FailNotEqual([]byte(expected), actual)
}

func assertServerResponseErrorCode(req clientRequest, expected codes.Code) error {
	_, err := req.Do()
	if err == nil {
		if expected != codes.OK {
			return fmt.Errorf("got no error, want %q", expected) // nolint: goerr113
		}

		return nil
	}

	actual := status.Convert(err).Code()

	if expected != actual {
		return fmt.Errorf("unexpected error code, got %q, want %q", actual, expected) // nolint: goerr113
	}

	return nil
}

func assertServerResponseErrorMessage(req clientRequest, expected string) error {
	_, err := req.Do()
	if err == nil {
		if expected != "" {
			return fmt.Errorf("got no error, want %q", expected) // nolint: goerr113
		}

		return nil
	}

	actual := status.Convert(err).Message()

	if expected != actual {
		return fmt.Errorf("unexpected error message, got %q, want %q", actual, expected) // nolint: goerr113
	}

	return nil
}
