export async function checkout(input: CheckoutInput) { return fetch('/api/orders', {method: 'POST', body: JSON.stringify(input)}); }
type CheckoutInput = { sku: string; quantity: number };
