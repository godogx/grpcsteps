Feature: Create Items with method unimplemented error

    Scenario: Method is unimplemented
        When I request a gRPC method "/grpctest.ItemService/CreateItems" with payload:
        """
        []
        """

        Then I should have a gRPC response with code "Unimplemented" and error message "CreateItems is not implemented"
