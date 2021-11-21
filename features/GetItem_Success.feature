Feature: Get Item (Success)

    Scenario: Get item without locale
        When I request a gRPC method "/grpctest.ItemService/GetItem" with payload:
        """
        {
            "id": 42
        }
        """

        Then I should have a gRPC response with code "OK"
        And I should have a gRPC response with payload:
        """
        {
            "id": 42,
            "name": "Test",
            "create_time": "<ignore-diff>"
        }
        """

    Scenario: Get item with locale
        When I request a gRPC method "/grpctest.ItemService/GetItem" with payload:
        """
        {
            "id": 42
        }
        """
        And The gRPC request has a header "Locale: en-US"

        Then I should have a gRPC response with code "OK"
        And I should have a gRPC response with payload:
        """
        {
            "id": 42,
            "locale": "en-US",
            "name": "Test",
            "create_time": "<ignore-diff>"
        }
        """
