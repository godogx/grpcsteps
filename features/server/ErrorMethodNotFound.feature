Feature: Error when method is not found

    Scenario: method not found
        Given "item-service" receives a grpc request "not-found"
