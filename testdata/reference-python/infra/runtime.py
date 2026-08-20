import os

def database_path() -> str:
    return os.environ.get("ORDER_DB", "orders.ndjson")
