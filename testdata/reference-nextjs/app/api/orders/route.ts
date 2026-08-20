import { createOrder } from "../../../src/order-service";

export async function POST(request: Request): Promise<Response> {
  const input = await request.json();
  return Response.json(await createOrder(input.productId));
}
