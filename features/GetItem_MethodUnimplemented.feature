Feature: Get Item

    Scenario: Method is unimplemented
        When I request a GRPC method "/grpctest.ItemService/GetItem" with payload:
        """
        {
            "id": 42
        }
        """

        Then I should have a GRPC response with code "Unimplemented" and error message "GetItem is not implemented"
