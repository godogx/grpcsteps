Feature: Error when service is not found

    Scenario: service is not registered
        Given "not-found" receives a grpc request "method"
