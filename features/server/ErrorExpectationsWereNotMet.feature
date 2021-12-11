Feature: Error when expectations were not met

    Scenario: there is one expectation
        Given "item-service" receives a grpc request "/grpctest.ItemService/GetItem" with payload:
        """
        {
            "id": 42
        }
        """
        And the grpc service responds with payload:
        """
        {
            "id": 42
            "name": "Item #42"
        }
        """
