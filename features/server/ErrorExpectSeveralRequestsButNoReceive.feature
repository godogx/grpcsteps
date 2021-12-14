Feature: Expect several requests but receive nothing

    Scenario: expect several requests
        Given "item-service" receives several grpc requests "/grpctest.ItemService/GetItem" with payload:
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
