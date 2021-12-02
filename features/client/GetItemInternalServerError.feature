Feature: Get Item with internal server error

    Scenario: Internal Server Error
        When I request a gRPC method "/grpctest.ItemService/GetItem" with payload:
        """
        {
            "id": 42
        }
        """

        Then I should have a gRPC response with code "INTERNAL"
        And I should have a gRPC response with error message "Internal Server Error"
