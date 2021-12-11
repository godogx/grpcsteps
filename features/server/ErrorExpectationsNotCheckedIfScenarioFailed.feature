Feature: Expectations are not checked when scenario is failed

    Scenario: there is one expectation
        Given "item-service" receives a grpc request "/grpctest.ItemService/GetItem" with payload:
        """
        {
            "id": 42
        }
        """
        And the grpc service responds with code "this fails"
