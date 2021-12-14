Feature: Expect a request but receive nothing

    Scenario: there is one request
        Given "item-service" receives a grpc request "/grpctest.ItemService/GetItem" with payload:
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
