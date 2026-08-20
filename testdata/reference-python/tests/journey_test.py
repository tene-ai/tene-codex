import json
import os
import tempfile
import unittest

from app.api import post_order
from app.worker import process_next

class OrderJourneyTest(unittest.TestCase):
    def test_api_queue_worker_database_flow(self):
        with tempfile.TemporaryDirectory() as root:
            database = os.path.join(root, "orders.ndjson")
            response = post_order({"product_id": "sku-42"})
            processed = process_next(database)
            with open(database, encoding="utf-8") as stored:
                persisted = json.loads(stored.read())
            self.assertEqual(response, {"id": "order-sku-42", "product_id": "sku-42"})
            self.assertEqual(processed, response)
            self.assertEqual(persisted, response)

if __name__ == "__main__":
    unittest.main()
