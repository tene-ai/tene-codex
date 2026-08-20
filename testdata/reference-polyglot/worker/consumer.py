def consume(message, queue, external_api):
    result = external_api.reserve(message["sku"])
    queue.ack(message)
    return result
