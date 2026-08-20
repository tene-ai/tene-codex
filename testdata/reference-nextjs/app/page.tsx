import { submitOrder } from "../src/order-client";

export async function CheckoutPage(productId: string): Promise<string> {
  const order = await submitOrder(productId);
  return `confirmed:${order.id}`;
}
