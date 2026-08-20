MESSAGES: list[dict] = []

def publish(message: dict) -> None:
    MESSAGES.append(message)

def receive() -> dict:
    return MESSAGES.pop(0)
