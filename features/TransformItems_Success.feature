Feature: Transform Items (Success)

    Scenario: Transform items
        When I request a gRPC method "/grpctest.ItemService/TransformItems" with payload:
        """
        [
            {
                "id": 42,
                "name": "Item #42"
            },
            {
                "id": 43,
                "name": "Item #42"
            }
        ]
        """

        Then I should have a gRPC response with code "OK"
        And I should have a gRPC response with payload:
        """
        [
            {
                "id": 42,
                "name": "Modified Item #42"
            },
            {
                "id": 43,
                "name": "Modified Item #42"
            }
        ]
        """
