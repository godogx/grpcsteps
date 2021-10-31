# Cucumber gRPC steps for Golang

[![GitHub Releases](https://img.shields.io/github/v/release/godogx/grpcsteps)](https://github.com/godogx/grpcsteps/releases/latest)
[![Build Status](https://github.com/godogx/grpcsteps/actions/workflows/test.yaml/badge.svg)](https://github.com/godogx/grpcsteps/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/godogx/grpcsteps/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/godogx/grpcsteps)
[![Go Report Card](https://goreportcard.com/badge/github.com/nhatthm/httpmock)](https://goreportcard.com/report/github.com/nhatthm/httpmock)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/github.com/godogx/grpcsteps)

`grpcsteps` uses [`nhatthm/grpcmock`](https://github.com/nhatthm/grpcmock) to provide steps for [`cucumber/godog`](https://github.com/cucumber/godog) and makes
it easy to run tests with grpc server and client.

Read more about [`nhatthm/grpcmock`](https://github.com/nhatthm/grpcmock)

## Prerequisites

- `Go >= 1.17`

## Usage

### Test a gPRC Server.

Initiate a client and register it to the scenario.

```go
package mypackage

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/cucumber/godog"
	"google.golang.org/grpc"

	"github.com/godogx/grpcsteps"
)

func TestIntegration(t *testing.T) {
	out := bytes.NewBuffer(nil)

	// Create a new grpc client.
	c := grpcsteps.NewClient(
		grpcsteps.RegisterService(
			grpctest.RegisterItemServiceServer,
			grpcsteps.WithDialOptions(
				grpc.WithInsecure(),
			),
		),
	)

	suite := godog.TestSuite{
		Name:                 "Integration",
		TestSuiteInitializer: nil,
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			// Register the client.
			c.RegisterContext(ctx)
		},
		Options: &godog.Options{
			Strict:    true,
			Output:    out,
			Randomize: rand.Int63(),
		},
	}

	// Run the suite.
	if status := suite.Run(); status != 0 {
		t.Fatal(out.String())
	}
}
```

#### Setup

In order to test the gPRC server, you have to register it to the client with `grpcsteps.RegisterService()` while initializing. The first argument is the
function that prototool generates for you. something like this

```go
package mypackage

func RegisterItemServiceServer(s grpc.ServiceRegistrar, srv ItemServiceServer) {
	s.RegisterService(&ItemService_ServiceDesc, srv)
}
```

You can configure how the client connects to the server by putting the options. For example:

```go
package mypackage

func createClient() *grpcsteps.Client {
	return grpcsteps.NewClient(
		grpcsteps.RegisterService(
			grpctest.RegisterItemServiceServer,
			grpcsteps.WithDialOptions(
				grpc.WithInsecure(),
			),
		),
	)
}
```

If you have multiple services and want to apply a same set of options to all, use `grpcsteps.WithDefaultServiceOptions()`. For example:

```go
package mypackage

func createClient() *grpcsteps.Client {
	return grpcsteps.NewClient(
		// Set default service options.
		grpcsteps.WithDefaultServiceOptions(
			grpcsteps.WithDialOptions(
                grpc.WithInsecure(),
            ),
		),
		// Register other services after this.
		grpcsteps.RegisterService(grpctest.RegisterItemServiceServer),
	)
}
```

#### Options

The options are:

- `grpcsteps.WithAddressProvider(interface{Addr() net.Addr})`: Connect to the server using the given address provider, the golang's built-in `*net.Listener` is
  an address provider.
- `grpcsteps.WithAddr(string)`: Connect to the server using the given address. For example: `:9090` or `localhost:9090`.
- `grpcsteps.WithDialOption(grpc.DialOption)`: Add a dial option for connecting to the server.
- `grpcsteps.WithDialOptions(...grpc.DialOption)`: Add multiple dial options for connecting to the server.

#### Steps

##### Prepare for a request

Create a new request with this pattern <br/>
`^I request(?: a)? (?:GRPC|grpc)(?: method)? "([^"]*)" with payload:?$`

Optionally, you can:

- Add a header to the request with <br/>
  `^The (?:GRPC|grpc) request has(?: a)? header "([^"]*): ([^"]*)"$`
- Set a timeout for the request with <br/>
  `^The (?:GRPC|grpc) request timeout is "([^"]*)"$`

For example:

```gherkin
Feature: Get Item

    Scenario: Get item with locale
        When I request a GRPC method "/grpctest.ItemService/GetItem" with payload:
        """
        {
            "id": 42
        }
        """
        And The GRPC request has a header "Locale: en-US"
```

##### Execute the request and validate the result.

- Check only the response code <br/>
  `^I should have(?: a)? (?:GRPC|grpc) response with code "([^"]*)"$`
- Check if the request is successful and the response payload matches an expectation <br/>
  `^I should have(?: a)? (?:GRPC|grpc) response with payload:?$`
- Check for error code and error message <br/>
  `^I should have(?: a)? (?:GRPC|grpc) response with error (?:message )?"([^"]*)"$` <br/>
  `^I should have(?: a)? (?:GRPC|grpc) response with code "([^"]*)" and error (?:message )?"([^"]*)"$`

For example:

```gherkin
Feature: Create Items

    Scenario: Create items
        When I request a GRPC method "/grpctest.ItemService/CreateItems" with payload:
        """
        [
            {
                "id": 42,
                "name": "Item #42"
            },
            {
                "id": 43,
                "name": "Item #42"
            }
        ]
        """

        Then I should have a GRPC response with payload:
        """
        {
            "num_items": 2
        }
        """
```
