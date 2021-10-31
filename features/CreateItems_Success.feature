Feature: Create Items (Success)

    Scenario: Create items
        When I request a GRPC method "/grpctest.ItemService/CreateItems" with payload:
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
        {
            "num_items": 2
        }
        """
