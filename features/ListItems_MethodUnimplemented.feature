Feature: List Items

    Scenario: Method is unimplemented
        When I request a GRPC method "/grpctest.ItemService/ListItems" with payload:
        """
        {}
        """

        Then I should have a GRPC response with code "Unimplemented" and error message "ListItems is not implemented"
