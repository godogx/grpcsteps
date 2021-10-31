Feature: Transform Items (Success)

    Scenario: Transform items
        When I request a GRPC method "/grpctest.ItemService/TransformItems" with payload:
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

        Then I should have a GRPC response with code "OK"
        And I should have a GRPC response with payload:
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
