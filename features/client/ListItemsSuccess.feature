Feature: List Items (Success)

    Scenario: List items without locale
        When I request a gRPC method "/grpctest.ItemService/ListItems" with payload:
        """
        {}
        """

        Then I should have a gRPC response with code "OK"
        And I should have a gRPC response with payload:
        """
        [
            {
                "id": 42,
                "name": "Test",
                "create_time": "<ignore-diff>"
            }
        ]
        """

    Scenario: List items with locale
        When I request a gRPC method "/grpctest.ItemService/ListItems" with payload:
        """
        {}
        """
        And The gRPC request has a header "Locale: en-US"

        Then I should have a gRPC response with code "OK"
        And I should have a gRPC response with payload:
        """
        [
            {
                "id": 42,
                "locale": "en-US",
                "name": "Test",
                "create_time": "<ignore-diff>"
            }
        ]
        """
