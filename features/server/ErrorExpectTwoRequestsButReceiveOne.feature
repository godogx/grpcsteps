Feature: Expect 2 requests but receive only 1

    Scenario: receive only 1 request
        Given "item-service" receives 2 grpc requests "/grpctest.ItemService/GetItem" with payload:
        """
        {
            "id": 42
        }
        """
        And the grpc service responds with payload:
        """
        {
            "id": 42,
            "name": "Item #42"
        }
        """

        When I request a grpc method "/grpctest.ItemService/GetItem" with payload:
        """
        {
            "id": 42
        }
        """

        Then I should have a grpc response with payload:
        """
        {
            "id": 42,
            "name": "Item #42"
        }
        """
