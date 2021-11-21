Feature: Create Items (Success)

    Scenario: Create items
        When I request a gRPC method "/grpctest.ItemService/CreateItems" with payload:
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
        {
            "num_items": 2
        }
        """
