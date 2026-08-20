from app.queue import receive
from app.repository import save

def process_next(database: str) -> dict:
    return save(database, receive())
