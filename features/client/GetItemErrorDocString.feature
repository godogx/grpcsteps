Feature: Get Item with error in doc string

    Scenario: With only error message
        When I request a gRPC method "/grpctest.ItemService/GetItem" with payload:
        """
        {
            "id": 42
        }
        """

        Then I should have a gRPC response with error message:
        """
        invalid "id"
        """

    Scenario: With code and error message
        When I request a gRPC method "/grpctest.ItemService/GetItem" with payload:
        """
        {
            "id": 42
        }
        """

        Then I should have a gRPC response with code "FailedPrecondition" and error message:
        """
        invalid "id"
        """
