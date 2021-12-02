Feature: Get Item with method unimplemented error

    Scenario: Method is unimplemented
        When I request a gRPC method "/grpctest.ItemService/GetItem" with payload:
        """
        {
            "id": 42
        }
        """

        Then I should have a gRPC response with code "Unimplemented" and error message "GetItem is not implemented"
