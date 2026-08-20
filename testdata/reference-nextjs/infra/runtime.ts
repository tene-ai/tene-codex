export function orderDatabasePath(environment: string): string {
  return environment === "test" ? "orders.test.ndjson" : "orders.ndjson";
}
