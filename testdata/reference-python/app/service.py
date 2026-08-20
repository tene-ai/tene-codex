from app.queue import publish

def accept_order(product_id: str) -> dict:
    order = {"id": "order-" + product_id, "product_id": product_id}
    publish(order)
    return order
