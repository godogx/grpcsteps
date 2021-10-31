Feature: Get Item

    Scenario: Item not found
        When I request a GRPC method "/grpctest.ItemService/GetItem" with payload:
        """
        {
            "id": 42
        }
        """

        Then I should have a GRPC response with code "NOT_FOUND"
        Then I should have a GRPC response with error message "Item 42 not found"
