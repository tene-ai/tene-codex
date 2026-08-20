import json

def save(database: str, order: dict) -> dict:
    with open(database, "a", encoding="utf-8") as output:
        output.write(json.dumps(order) + "\n")
    return order
