Feature:
    Scenario: Server is not online
        When I request a GRPC method "/NoServer/GetItem" with payload:
        """
        {
            "id": 42
        }
        """
        And The GRPC request timeout is "1ms"

        Then I should have a GRPC response with code "DEADLINE_EXCEEDED"
