import { appendFile } from "node:fs/promises";

export async function saveOrder(order: {id: string; productId: string}): Promise<{id: string}> {
  await appendFile(process.env.ORDER_DB!, JSON.stringify(order) + "\n");
  return order;
}
