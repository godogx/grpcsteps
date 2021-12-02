Feature: Server is not ready

    Scenario: Server is not online
        When I request a gRPC method "/NoServer/GetItem" with payload:
        """
        {
            "id": 42
        }
        """
        And The gRPC request timeout is "1ms"

#        Then I should have a gRPC response with code "DEADLINE_EXCEEDED"

        Then I should have a gRPC response with error message:
        """
        context deadline exceeded
        """
