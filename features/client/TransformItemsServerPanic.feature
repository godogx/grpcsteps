Feature: Transform Items

    Scenario: Server Panic
        When I request a gRPC method "/grpctest.ItemService/TransformItems" with payload:
        """
        []
        """

        Then I should have a gRPC response with code "Internal"
