Feature: Create Items

    Scenario: Method is unimplemented
        When I request a GRPC method "/grpctest.ItemService/CreateItems" with payload:
        """
        []
        """

        Then I should have a GRPC response with code "Unimplemented" and error message "CreateItems is not implemented"
