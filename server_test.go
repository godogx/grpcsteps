package grpcsteps_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/nhatthm/grpcmock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/godogx/grpcsteps"
	"github.com/godogx/grpcsteps/internal/grpctest"
)

func TestExternalServiceManager_Success(t *testing.T) {
	t.Parallel()

	runServerTest(t, "Success")
}

func TestExternalServiceManager_Error(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		expected string
	}{
		{
			scenario: "ErrorServiceNotFound",
			expected: `Feature: Error when service is not found

  Scenario: service is not registered                  # features/server/ErrorServiceNotFound.feature:3
    Given "not-found" receives a grpc request "method" # server.go:86 -> *ExternalServiceManager
    grpc service not found, did you forget to setup the grpc service "not-found"?
`,
		},
		{
			scenario: "ErrorMethodNotFound",
			expected: `Feature: Error when method is not found

  Scenario: method not found                                 # features/server/ErrorMethodNotFound.feature:3
    Given "item-service" receives a grpc request "not-found" # server.go:86 -> *ExternalServiceManager
    grpc method not found: not-found
`,
		},
		{
			scenario: "ErrorMethodNotSupported",
			expected: `Feature: Error when method is not supported

  Scenario: method not supported                                                        # features/server/ErrorMethodNotSupported.feature:3
    Given "item-service" receives a grpc request "/grpctest.ItemService/TransformItems" # server.go:86 -> *ExternalServiceManager
    grpc method not supported: BidirectionalStream /grpctest.ItemService/TransformItems
`,
		},
		{
			scenario: "ErrorExpectationsWereNotMet",
			expected: `Feature: Error when expectations were not met

  Scenario: there is one expectation                                                           # features/server/ErrorExpectationsWereNotMet.feature:3
    Given "item-service" receives a grpc request "/grpctest.ItemService/GetItem" with payload: # server.go:94 -> *ExternalServiceManager
      """
      {
          "id": 42
      }
      """
    And the grpc service responds with payload:                                                # server.go:115 -> *ExternalServiceManager
      """
      {
          "id": 42
          "name": "Item #42"
      }
      """
    after scenario hook failed: there are remaining expectations that were not met:
- Unary /grpctest.ItemService/GetItem
    with payload using matcher.JSONMatcher
        {
    "id": 42
}
`,
		},
		{
			scenario: "ErrorExpectationsNotCheckedIfScenarioFailed",
			expected: `Feature: Expectations are not checked when scenario is failed

  Scenario: there is one expectation                                                           # features/server/ErrorExpectationsNotCheckedIfScenarioFailed.feature:3
    Given "item-service" receives a grpc request "/grpctest.ItemService/GetItem" with payload: # server.go:94 -> *ExternalServiceManager
      """
      {
          "id": 42
      }
      """
    And the grpc service responds with code "this fails"                                       # server.go:145 -> *ExternalServiceManager
    invalid code: "\"THIS FAILS\""

--- Failed steps:

  Scenario: there is one expectation # features/server/ErrorExpectationsNotCheckedIfScenarioFailed.feature:3
    And the grpc service responds with code "this fails" # features/server/ErrorExpectationsNotCheckedIfScenarioFailed.feature:10
      Error: invalid code: "\"THIS FAILS\""
`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			tt := &testT{}

			runServerTest(tt, tc.scenario, noColors())

			assert.Error(t, tt.error)
			assert.NotEmpty(t, tc.expected)
			assert.Contains(t, tt.error.Error(), tc.expected)
		})
	}
}

func runServerTest(
	t suiteT,
	scenario string,
	opts ...suiteOption,
) {
	buf := bufconn.Listen(1024 * 1024)

	srv := grpcsteps.NewExternalServiceManager()
	srv.AddService("item-service",
		grpcmock.WithListener(buf),
		grpcmock.RegisterService(grpctest.RegisterItemServiceServer),
	)

	c := grpcsteps.NewClient(
		grpcsteps.WithDefaultServiceOptions(
			grpcsteps.WithDialOptions(
				grpc.WithInsecure(),
				grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
					return buf.Dial()
				}),
			),
		),
		grpcsteps.RegisterService(grpctest.RegisterItemServiceServer),
	)

	opts = append(opts,
		afterSuite(func() {
			srv.Close()
		}),
		initScenario(c.RegisterContext, srv.RegisterContext),
		featureFiles(fmt.Sprintf("features/server/%s.feature", scenario)),
	)

	runSuite(t, opts...)
}
