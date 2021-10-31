Feature: Transform Items

    Scenario: Method is unimplemented
        When I request a GRPC method "/grpctest.ItemService/TransformItems" with payload:
        """
        []
        """

        Then I should have a GRPC response with code "Unimplemented" and error message "TransformItems is not implemented"
