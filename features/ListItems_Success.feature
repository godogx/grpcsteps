Feature: List Items (Success)

    Scenario: List items without locale
        When I request a GRPC method "/grpctest.ItemService/ListItems" with payload:
        """
        {}
        """

        Then I should have a GRPC response with code "OK"
        And I should have a GRPC response with payload:
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
        When I request a GRPC method "/grpctest.ItemService/ListItems" with payload:
        """
        {}
        """
        And The GRPC request has a header "Locale: en-US"

        Then I should have a GRPC response with code "OK"
        And I should have a GRPC response with payload:
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
