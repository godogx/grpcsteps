Feature: Get Item

    Scenario: Server Panic
        When I request a GRPC method "/grpctest.ItemService/GetItem" with payload:
        """
        {
            "id": 42
        }
        """

        Then I should have a GRPC response with code "Internal"
