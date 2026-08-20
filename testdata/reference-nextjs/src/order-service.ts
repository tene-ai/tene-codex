import { saveOrder } from "./order-repository";

export async function createOrder(productId: string): Promise<{id: string}> {
  return saveOrder({id: `order-${productId}`, productId});
}
