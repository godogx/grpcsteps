Feature: Create Items

    Scenario: With only error message
        When I request a GRPC method "/grpctest.ItemService/CreateItems" with payload:
        """
        []
        """

        Then I should have a GRPC response with error message:
        """
        invalid "name"
        """

    Scenario: With code and error message
        When I request a GRPC method "/grpctest.ItemService/CreateItems" with payload:
        """
        []
        """

        Then I should have a GRPC response with code "FailedPrecondition" and error message:
        """
        invalid "name"
        """
