Feature: Error when method is not supported

    Scenario: method not supported
        Given "item-service" receives a grpc request "/grpctest.ItemService/TransformItems"
