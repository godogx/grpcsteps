Feature: List Items

    Scenario: With only error message
        When I request a GRPC method "/grpctest.ItemService/ListItems" with payload:
        """
        {}
        """

        Then I should have a GRPC response with error message:
        """
        invalid "page_size"
        """

    Scenario: With code and error message
        When I request a GRPC method "/grpctest.ItemService/ListItems" with payload:
        """
        {}
        """

        Then I should have a GRPC response with code "FailedPrecondition" and error message:
        """
        invalid "page_size"
        """
