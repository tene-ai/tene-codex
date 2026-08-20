from app.service import accept_order

def post_order(payload: dict) -> dict:
    return accept_order(payload["product_id"])
