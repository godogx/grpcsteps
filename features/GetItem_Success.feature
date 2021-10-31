Feature: Get Item (Success)

    Scenario: Get item without locale
        When I request a GRPC method "/grpctest.ItemService/GetItem" with payload:
        """
        {
            "id": 42
        }
        """

        Then I should have a GRPC response with code "OK"
        And I should have a GRPC response with payload:
        """
        {
            "id": 42,
            "name": "Test",
            "create_time": "<ignore-diff>"
        }
        """

    Scenario: Get item with locale
        When I request a GRPC method "/grpctest.ItemService/GetItem" with payload:
        """
        {
            "id": 42
        }
        """
        And The GRPC request has a header "Locale: en-US"

        Then I should have a GRPC response with code "OK"
        And I should have a GRPC response with payload:
        """
        {
            "id": 42,
            "locale": "en-US",
            "name": "Test",
            "create_time": "<ignore-diff>"
        }
        """
