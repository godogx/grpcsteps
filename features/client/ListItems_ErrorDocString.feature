Feature: List Items

    Scenario: With only error message
        When I request a gRPC method "/grpctest.ItemService/ListItems" with payload:
        """
        {}
        """

        Then I should have a gRPC response with error message:
        """
        invalid "page_size"
        """

    Scenario: With code and error message
        When I request a gRPC method "/grpctest.ItemService/ListItems" with payload:
        """
        {}
        """

        Then I should have a gRPC response with code "FailedPrecondition" and error message:
        """
        invalid "page_size"
        """
