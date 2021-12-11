Feature: Test all supported types with all statements

    Scenario Outline: Return code
        Given "item-service" receives a grpc request "/grpctest.ItemService/<method>"
        And the grpc service responds with code "InvalidArgument"

        When I request a grpc method "/grpctest.ItemService/<method>" with payload:
        """
        <request>
        """

        Then I should have a grpc response with code "InvalidArgument"

        Examples:
            | method      | request      |
            | GetItem     | {"id": 42}   |
            | ListItems   | {}           |
            | CreateItems | [{"id": 42}] |

    Scenario Outline: Return error message
        Given "item-service" receives a grpc request "/grpctest.ItemService/<method>"
        And the grpc service responds with error "Internal Server Error"

        When I request a grpc method "/grpctest.ItemService/<method>" with payload:
        """
        <request>
        """

        Then I should have a grpc response with error message "Internal Server Error"

        Examples:
            | method      | request      |
            | GetItem     | {"id": 42}   |
            | ListItems   | {}           |
            | CreateItems | [{"id": 42}] |

    Scenario Outline: Return error message in doc string
        Given "item-service" receives a grpc request "/grpctest.ItemService/<method>"
        And the grpc service responds with error:
        """
        Internal Server Error
        """

        When I request a grpc method "/grpctest.ItemService/<method>" with payload:
        """
        <request>
        """

        Then I should have a grpc response with error message "Internal Server Error"

        Examples:
            | method      | request      |
            | GetItem     | {"id": 42}   |
            | ListItems   | {}           |
            | CreateItems | [{"id": 42}] |

    Scenario Outline: Return error code and message
        Given "item-service" receives a grpc request "/grpctest.ItemService/<method>"
        And the grpc service responds with code "Internal" and error "Internal Server Error"

        When I request a grpc method "/grpctest.ItemService/<method>" with payload:
        """
        <request>
        """

        Then I should have a grpc response with code "Internal" and error message "Internal Server Error"

        Examples:
            | method      | request      |
            | GetItem     | {"id": 42}   |
            | ListItems   | {}           |
            | CreateItems | [{"id": 42}] |

    Scenario Outline: Return error code and message in doc string
        Given "item-service" receives a grpc request "/grpctest.ItemService/<method>"
        And the grpc service responds with code "Internal" and error:
        """
        Internal Server Error
        """

        When I request a grpc method "/grpctest.ItemService/<method>" with payload:
        """
        <request>
        """

        Then I should have a grpc response with code "Internal" and error message "Internal Server Error"

        Examples:
            | method      | request      |
            | GetItem     | {"id": 42}   |
            | ListItems   | {}           |
            | CreateItems | [{"id": 42}] |

    Scenario Outline: With the same header
        Given "item-service" receives a grpc request "/grpctest.ItemService/<method>"
        And the grpc request has a header "Locale: en-US"
        And the grpc service responds with code "InvalidArgument"

        When I request a grpc method "/grpctest.ItemService/<method>" with payload:
        """
        <request>
        """

        Then I should have a grpc response with code "Internal" and error:
        """
        Expected: <type> /grpctest.ItemService/<method>
            with header:
                Locale: en-US
        Actual: <type> /grpctest.ItemService/<method>
            with payload
                <request>
        Error: header "Locale" with value "en-US" expected, "" received

        """

        When I request a grpc method "/grpctest.ItemService/<method>" with payload:
        """
        <request>
        """
        And the grpc request has a header "Locale: en-US"

        Then I should have a grpc response with code "InvalidArgument"

        Examples:
            | method      | type         | request     |
            | GetItem     | Unary        | {"id":42}   |
            | ListItems   | ServerStream | {}          |
            | CreateItems | ClientStream | [{"id":42}] |

    Scenario Outline: With payload
        Given "item-service" receives a grpc request "/grpctest.ItemService/<method>" with payload:
        """
        <expect>
        """
        And the grpc service responds with code "InvalidArgument"

        When I request a grpc method "/grpctest.ItemService/<method>" with payload:
        """
        <request>
        """

        Then I should have a grpc response with code "InvalidArgument"

        Examples:
            | method      | expect                    | request     |
            | GetItem     | {"id": "<ignore-diff>"}   | {"id":42}   |
            | ListItems   | {}                        | {}          |
            | CreateItems | [{"id": "<ignore-diff>"}] | [{"id":42}] |

    Scenario Outline: With payload from file in doc string
        Given "item-service" receives a grpc request "/grpctest.ItemService/<method>" with payload from file:
        """
        <file>
        """
        And the grpc service responds with code "InvalidArgument"

        When I request a grpc method "/grpctest.ItemService/<method>" with payload:
        """
        <request>
        """

        Then I should have a grpc response with code "InvalidArgument"

        Examples:
            | method      | file                                        | request     |
            | GetItem     | resources/fixtures/expect-get-item.json     | {"id":42}   |
            | ListItems   | resources/fixtures/expect-list-items.json   | {}          |
            | CreateItems | resources/fixtures/expect-create-items.json | [{"id":42}] |

    Scenario Outline: With payload from file
        Given "item-service" receives a grpc request "/grpctest.ItemService/<method>" with payload from file "<file>"
        And the grpc service responds with code "InvalidArgument"

        When I request a grpc method "/grpctest.ItemService/<method>" with payload:
        """
        <request>
        """

        Then I should have a grpc response with code "InvalidArgument"

        Examples:
            | method      | file                                        | request     |
            | GetItem     | resources/fixtures/expect-get-item.json     | {"id":42}   |
            | ListItems   | resources/fixtures/expect-list-items.json   | {}          |
            | CreateItems | resources/fixtures/expect-create-items.json | [{"id":42}] |

    Scenario Outline: Response from doc string
        Given "item-service" receives a grpc request "/grpctest.ItemService/<method>" with payload from file "<request_file>"
        And the grpc service responds with payload:
        """
        <response>
        """

        When I request a grpc method "/grpctest.ItemService/<method>" with payload from file:
        """
        <request_file>
        """

        Then I should have a grpc response with payload:
        """
        <response>
        """

        Examples:
            | method      | request_file                                 | response        |
            | GetItem     | resources/fixtures/request-get-item.json     | {"id":42}       |
            | ListItems   | resources/fixtures/request-list-items.json   | [{"id":42}]     |
            | CreateItems | resources/fixtures/request-create-items.json | {"num_items":3} |

    Scenario Outline: Response from file
        Given "item-service" receives a grpc request "/grpctest.ItemService/<method>" with payload from file "<request_file>"
        And the grpc service responds with payload from file "<response_file>"

        When I request a grpc method "/grpctest.ItemService/<method>" with payload from file:
        """
        <request_file>
        """

        Then I should have a grpc response with payload from file:
        """
        <response_file>
        """

        Examples:
            | method      | request_file                                 | response_file                                 |
            | GetItem     | resources/fixtures/request-get-item.json     | resources/fixtures/response-get-item.json     |
            | ListItems   | resources/fixtures/request-list-items.json   | resources/fixtures/response-list-items.json   |
            | CreateItems | resources/fixtures/request-create-items.json | resources/fixtures/response-create-items.json |

    Scenario Outline: Response from file in doc string
        Given "item-service" receives a grpc request "/grpctest.ItemService/<method>" with payload from file "<request_file>"
        And the grpc service responds with payload from file:
        """
        <response_file>
        """

        When I request a grpc method "/grpctest.ItemService/<method>" with payload from file:
        """
        <request_file>
        """

        Then I should have a grpc response with payload from file:
        """
        <response_file>
        """

        Examples:
            | method      | request_file                                 | response_file                                 |
            | GetItem     | resources/fixtures/request-get-item.json     | resources/fixtures/response-get-item.json     |
            | ListItems   | resources/fixtures/request-list-items.json   | resources/fixtures/response-list-items.json   |
            | CreateItems | resources/fixtures/request-create-items.json | resources/fixtures/response-create-items.json |
