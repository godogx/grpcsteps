Feature: Get Item

    Scenario: With only error message
        When I request a GRPC method "/grpctest.ItemService/GetItem" with payload:
        """
        {
            "id": 42
        }
        """

        Then I should have a GRPC response with error message:
        """
        invalid "id"
        """

    Scenario: With code and error message
        When I request a GRPC method "/grpctest.ItemService/GetItem" with payload:
        """
        {
            "id": 42
        }
        """

        Then I should have a GRPC response with code "FailedPrecondition" and error message:
        """
        invalid "id"
        """
