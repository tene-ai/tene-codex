export async function submitOrder(productId: string): Promise<{id: string}> {
  const response = await fetch("/api/orders", {method: "POST", body: JSON.stringify({productId})});
  return response.json();
}
