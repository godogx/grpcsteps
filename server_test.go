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
    Given "not-found" receives a grpc request "method" # server.go:97 -> *ExternalServiceManager
    grpc service not found, did you forget to setup the grpc service "not-found"?
`,
		},
		{
			scenario: "ErrorMethodNotFound",
			expected: `Feature: Error when method is not found

  Scenario: method not found                                 # features/server/ErrorMethodNotFound.feature:3
    Given "item-service" receives a grpc request "not-found" # server.go:97 -> *ExternalServiceManager
    grpc method not found: not-found
`,
		},
		{
			scenario: "ErrorMethodNotSupported",
			expected: `Feature: Error when method is not supported

  Scenario: method not supported                                                        # features/server/ErrorMethodNotSupported.feature:3
    Given "item-service" receives a grpc request "/grpctest.ItemService/TransformItems" # server.go:97 -> *ExternalServiceManager
    grpc method not supported: BidirectionalStream /grpctest.ItemService/TransformItems
`,
		},
		{
			scenario: "ErrorExpectOneRequestButNoReceive",
			expected: `Feature: Expect a request but receive nothing

  Scenario: there is one request                                                               # features/server/ErrorExpectOneRequestButNoReceive.feature:3
    Given "item-service" receives a grpc request "/grpctest.ItemService/GetItem" with payload: # server.go:105 -> *ExternalServiceManager
      """
      {
          "id": 42
      }
      """
    And the grpc service responds with payload:                                                # server.go:176 -> *ExternalServiceManager
      """
      {
          "id": 42,
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
			scenario: "ErrorExpectTwoRequestsButReceiveOne",
			expected: `Feature: Expect 2 requests but receive only 1

  Scenario: receive only 1 request                                                              # features/server/ErrorExpectTwoRequestsButReceiveOne.feature:3
    Given "item-service" receives 2 grpc requests "/grpctest.ItemService/GetItem" with payload: # server.go:130 -> *ExternalServiceManager
      """
      {
          "id": 42
      }
      """
    And the grpc service responds with payload:                                                 # server.go:176 -> *ExternalServiceManager
      """
      {
          "id": 42,
          "name": "Item #42"
      }
      """
    When I request a grpc method "/grpctest.ItemService/GetItem" with payload:                  # client.go:94 -> *Client
      """
      {
          "id": 42
      }
      """
    Then I should have a grpc response with payload:                                            # client.go:115 -> *Client
      """
      {
          "id": 42,
          "name": "Item #42"
      }
      """
    after scenario hook failed: there are remaining expectations that were not met:
- Unary /grpctest.ItemService/GetItem (called: 1 time(s), remaining: 1 time(s))
    with payload using matcher.JSONMatcher
        {
    "id": 42
}
`,
		},
		{
			scenario: "ErrorExpectSeveralRequestsButNoReceive",
			expected: `Feature: Expect several requests but receive nothing

  Scenario: expect several requests                                                                   # features/server/ErrorExpectSeveralRequestsButNoReceive.feature:3
    Given "item-service" receives several grpc requests "/grpctest.ItemService/GetItem" with payload: # server.go:155 -> *ExternalServiceManager
      """
      {
          "id": 42
      }
      """
    And the grpc service responds with payload:                                                       # server.go:176 -> *ExternalServiceManager
      """
      {
          "id": 42,
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
    Given "item-service" receives a grpc request "/grpctest.ItemService/GetItem" with payload: # server.go:105 -> *ExternalServiceManager
      """
      {
          "id": 42
      }
      """
    And the grpc service responds with code "this fails"                                       # server.go:206 -> *ExternalServiceManager
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
