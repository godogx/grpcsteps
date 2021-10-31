Feature: Get Item

    Scenario: Internal Server Error
        When I request a GRPC method "/grpctest.ItemService/GetItem" with payload:
        """
        {
            "id": 42
        }
        """

        Then I should have a GRPC response with code "INTERNAL"
        Then I should have a GRPC response with error message "Internal Server Error"
