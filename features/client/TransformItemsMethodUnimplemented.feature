Feature: Transform Items

    Scenario: Method is unimplemented
        When I request a gRPC method "/grpctest.ItemService/TransformItems" with payload:
        """
        []
        """

        Then I should have a gRPC response with code "Unimplemented" and error message "TransformItems is not implemented"
