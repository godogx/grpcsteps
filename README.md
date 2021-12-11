# Cucumber gRPC steps for Golang

[![GitHub Releases](https://img.shields.io/github/v/release/godogx/grpcsteps)](https://github.com/godogx/grpcsteps/releases/latest)
[![Build Status](https://github.com/godogx/grpcsteps/actions/workflows/test.yaml/badge.svg)](https://github.com/godogx/grpcsteps/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/godogx/grpcsteps/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/godogx/grpcsteps)
[![Go Report Card](https://goreportcard.com/badge/github.com/nhatthm/httpmock)](https://goreportcard.com/report/github.com/nhatthm/httpmock)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/github.com/godogx/grpcsteps)

`grpcsteps` uses [`nhatthm/grpcmock`](https://github.com/nhatthm/grpcmock) to provide steps for [`cucumber/godog`](https://github.com/cucumber/godog) and makes
it easy to run tests with grpc server and client.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Usage](#usage)
    - [Mock external gPRC Services](#mock-external-gprc-services)
        - [Setup](#setup)
        - [Steps](#steps)
            - [Prepare for a request](#prepare-for-a-request)
            - [Response](#response)
    - [Test a gPRC Server](#test-a-gprc-server)
        - [Setup](#setup-1)
        - [Options](#options)
        - [Steps](#steps-1)
            - [Prepare for a request](#prepare-for-a-request-1)
            - [Execute the request and validate the result](#execute-the-request-and-validate-the-result)

## Prerequisites

- `Go >= 1.17`

[<sub><sup>[table of contents]</sup></sub>](#table-of-contents)

## Usage

### Mock external gPRC Services

This is for describing behaviors of gRPC endpoints that are called by the app during test (e.g. 3rd party APIs). The mock creates an gRPC server for each of
registered services and allows control of expected requests and responses with gherkin steps.

In simple case, you can define the expected method and response.

```gherkin
Feature: Get Item

    Scenario: Success
        Given "item-service" receives a grpc request "/grpctest.ItemService/GetItem" with payload:
        """
        {
            "id": 42
        }
        """

        And the grpc service responds with payload:
        """
        {
            "id": 42,
            "name": "Item #42"
        }
        """

        # Your application call.
```

For starting, initiate a server and register it to the scenario.

```go
package mypackage

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/cucumber/godog"

	"github.com/godogx/grpcsteps"
)

func TestIntegration(t *testing.T) {
	out := bytes.NewBuffer(nil)

	// Create a new grpc servers manager
	m := grpcsteps.NewExternalServiceManager()

	// Setup the 3rd party services here.

	suite := godog.TestSuite{
		Name: "Integration",
		TestSuiteInitializer: func(ctx *godog.TestSuiteContext) {
			ctx.After(func() {
				m.Close()
			})
		},
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			m.RegisterContext(ctx)
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

[<sub><sup>[table of contents]</sup></sub>](#table-of-contents)

#### Setup

In order to mock the gPRC server, you have to register it to the manager with `AddService()` while initializing. The first argument is the service ID, the
second argument is the function that prototool generates for you. Something like this:

```go
package mypackage

import "google.golang.org/grpc"

func RegisterItemServiceServer(s grpc.ServiceRegistrar, srv ItemServiceServer) {
	s.RegisterService(&ItemService_ServiceDesc, srv)
}
```

For example:

```go
package mypackage

import (
	"testing"

	"github.com/godogx/grpcsteps"
)

func TestIntegration(t *testing.T) {
	// Create a new grpc servers manager.
	m := grpcsteps.NewExternalServiceManager()

	itemServiceAddr := m.AddService("item-service", RegisterItemServiceServer)

	// itemServiceAddr is going to be something like "[::]:52299".
	// Use that addr for the client in the application.

	// Run test suite.
}
```

By default, the manager spins up a gRPC with a random port. If you don't like that, you can specify the one you like with `grpcmock.WithPort()`. For example:

```go
package mypackage

import (
	"testing"

	"github.com/nhatthm/grpcmock"

	"github.com/godogx/grpcsteps"
)

func TestIntegration(t *testing.T) {
	// Create a new grpc servers manager
	m := grpcsteps.NewExternalServiceManager()

	itemServiceAddr := m.AddService("item-service", RegisterItemServiceServer,
		grpcmock.WithPort(9000),
	)

	// itemServiceAddr is "[::]:9000".

	// Run test suite.
}
```

You can also use a listener, for example `bufconn`

```go
package mypackage

import (
	"testing"

	"github.com/nhatthm/grpcmock"
	"google.golang.org/grpc/test/bufconn"

	"github.com/godogx/grpcsteps"
)

func TestIntegration(t *testing.T) {
	buf := bufconn.Listen(1024 * 1024)

	// Create a new grpc servers manager
	m := grpcsteps.NewExternalServiceManager()

	m.AddService("item-service", RegisterItemServiceServer,
		grpcmock.WithListener(buf),
	)

	// In this case, use the `buf` to connect to server

	// Run test suite.
}
```

[<sub><sup>[table of contents]</sup></sub>](#table-of-contents)

#### Steps

##### Prepare for a request

Mock a new request with (one of) these patterns

- `^"([^"]*)" receives a (?:gRPC|GRPC|grpc) request "([^"]*)"$`
- `^"([^"]*)" receives a (?:gRPC|GRPC|grpc) request "([^"]*)" with payload:$`
- `^"([^"]*)" receives a (?:gRPC|GRPC|grpc) request "([^"]*)" with payload from file "([^"]+)"$`
- `^"([^"]*)" receives a (?:gRPC|GRPC|grpc) request "([^"]*)" with payload from file:$`

Optionally, you can:

- Add a header to the request with <br/>
  `^The (?:gRPC|GRPC|grpc) request has(?: a)? header "([^"]*): ([^"]*)"$`

For example:

```gherkin
Feature: Get Item

    Scenario: Get item with locale
        Given "item-service" receives a grpc request "/grpctest.ItemService/GetItem" with payload:
        """
        {
            "id": 42
        }
        """
        And The gRPC request has a header "Locale: en-US"

        # Your application call.
```

Note, you can use `"<ignore-diff>"` in the payload to tell the assertion to ignore a JSON field. For example:

```gherkin
Feature: Create Items

    Scenario: Create items
        Given "item-service" receives a grpc request "/grpctest.ItemService/CreateItems" with payload:
        """
        {
            "id": 42
            "name": "<ignore-diff>",
            "category": "<ignore-diff>",
            "metadata": "<ignore-diff>"
        }
        """

        And the gRPC service responds with payload:
        """
        {
            "num_items": 1
        }
        """

        # The application call.
        When my application receives a request to create some items:
        """
        [
            {
                "id": 42
                "name": "Item #42",
                "category": 40,
                "metadata": {
                    "tags": ["soup"]
                }
            }
        ]
        """

        Then 1 item is created
```

`"<ignore-diff>"` can ignore any types, not just string.

[<sub><sup>[table of contents]</sup></sub>](#table-of-contents)

##### Response

- Respond `OK` with payload <br/>
  `^[tT]he (?:gRPC|GRPC|grpc) service responds with payload:?$` <br/>
  `^[tT]he (?:gRPC|GRPC|grpc) service responds with payload from file "([^"]+)"$` <br/>
  `^[tT]he (?:gRPC|GRPC|grpc) service responds with payload from file:$`
- Response with code and error message <br/>
  `^[tT]he (?:gRPC|GRPC|grpc) service responds with code "([^"]*)"$` <br/>
  `^[tT]he (?:gRPC|GRPC|grpc) service responds with error (?:message )?"([^"]*)"$` <br/>
  `^[tT]he (?:gRPC|GRPC|grpc) service responds with code "([^"]*)" and error (?:message )?"([^"]*)"$` <br/>
  <br/>
  If your error message contains quotes `"`, better use these with a doc string<br/>
  `^[tT]he (?:gRPC|GRPC|grpc) service responds with error(?: message)?:$` </br>
  `^[tT]he (?:gRPC|GRPC|grpc) service responds with code "([^"]*)" and error(?: message)?:$` </br>

For example:

```gherkin
Feature: Create Items

    Scenario: Create items
        Given "item-service" receives a gRPC request "/grpctest.ItemService/CreateItems" with payload:
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

        And the gRPC service responds with payload:
        """
        {
            "num_items": 2
        }
        """

        # Your application call.
```

or

```gherkin
Feature: Create Items

    Scenario: Create items
        Given "item-service" receives a gRPC request "/grpctest.ItemService/CreateItems" with payload:
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

        And the gRPC service responds with code "InvalidArgument" and error "Invalid ID #42"
```

[<sub><sup>[table of contents]</sup></sub>](#table-of-contents)

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

[<sub><sup>[table of contents]</sup></sub>](#table-of-contents)

#### Setup

In order to test the gPRC server, you have to register it to the client with `grpcsteps.RegisterService()` while initializing. The first argument is the
function that prototool generates for you. Something like this:

```go
package mypackage

import "google.golang.org/grpc"

func RegisterItemServiceServer(s grpc.ServiceRegistrar, srv ItemServiceServer) {
	s.RegisterService(&ItemService_ServiceDesc, srv)
}
```

You can configure how the client connects to the server by putting the options. For example:

```go
package mypackage

import "google.golang.org/grpc"

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

import "google.golang.org/grpc"

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

[<sub><sup>[table of contents]</sup></sub>](#table-of-contents)

#### Options

The options are:

- `grpcsteps.WithAddressProvider(interface{Addr() net.Addr})`: Connect to the server using the given address provider, the golang's built-in `*net.Listener` is
  an address provider.
- `grpcsteps.WithAddr(string)`: Connect to the server using the given address. For example: `:9090` or `localhost:9090`.
- `grpcsteps.WithDialOption(grpc.DialOption)`: Add a dial option for connecting to the server.
- `grpcsteps.WithDialOptions(...grpc.DialOption)`: Add multiple dial options for connecting to the server.

[<sub><sup>[table of contents]</sup></sub>](#table-of-contents)

#### Steps

##### Prepare for a request

Create a new request with (one of) these patterns

- `^I request(?: a)? (?:gRPC|GRPC|grpc)(?: method)? "([^"]*)" with payload:?$`
- `^I request(?: a)? (?:gRPC|GRPC|grpc)(?: method)? "([^"]*)" with payload from file "([^"]+)"$`
- `^I request(?: a)? (?:gRPC|GRPC|grpc)(?: method)? "([^"]*)" with payload from file:$`

Optionally, you can:

- Add a header to the request with <br/>
  `^The (?:gRPC|GRPC|grpc) request has(?: a)? header "([^"]*): ([^"]*)"$`
- Set a timeout for the request with <br/>
  `^The (?:gRPC|GRPC|grpc) request timeout is "([^"]*)"$`

For example:

```gherkin
Feature: Get Item

    Scenario: Get item with locale
        When I request a gRPC method "/grpctest.ItemService/GetItem" with payload:
        """
        {
            "id": 42
        }
        """
        And The gRPC request has a header "Locale: en-US"
```

[<sub><sup>[table of contents]</sup></sub>](#table-of-contents)

##### Execute the request and validate the result.

- Check only the response code <br/>
  `^I should have(?: a)? (?:gRPC|GRPC|grpc) response with code "([^"]*)"$`
- Check if the request is successful and the response payload matches an expectation <br/>
  `^I should have(?: a)? (?:gRPC|GRPC|grpc) response with payload:?$`<br/>
  `^I should have(?: a)? (?:gRPC|GRPC|grpc) response with payload from file "([^"]+)"$`<br/>
  `^I should have(?: a)? (?:gRPC|GRPC|grpc) response with payload from file:$`
- Check for error code and error message <br/>
  `^I should have(?: a)? (?:gRPC|GRPC|grpc) response with error (?:message )?"([^"]*)"$` <br/>
  `^I should have(?: a)? (?:gRPC|GRPC|grpc) response with code "([^"]*)" and error (?:message )?"([^"]*)"$`<br/>
  <br/>
  If your error message contains quotes `"`, better use these with a doc string<br/>
  `^I should have(?: a)? (?:gRPC|GRPC|grpc) response with error (?:message )?:$` <br/>
  `^I should have(?: a)? (?:gRPC|GRPC|grpc) response with code "([^"]*)" and error (?:message )?:$`<br/>

For example:

```gherkin
Feature: Create Items

    Scenario: Create items
        When I request a gRPC method "/grpctest.ItemService/CreateItems" with payload:
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

        Then I should have a gRPC response with payload:
        """
        {
            "num_items": 2
        }
        """
```

or

```gherkin
Feature: Create Items

    Scenario: Create items
        When I request a gRPC method "/grpctest.ItemService/CreateItems" with payload:
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

        Then I should have a gRPC response with error:
        """
        invalid "id"
        """
```

[<sub><sup>[table of contents]</sup></sub>](#table-of-contents)
