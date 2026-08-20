import assert from "node:assert/strict";
import { mkdtemp, readFile, rm } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";

export async function runCheckoutJourney(productId, options = {}) {
  const root = await mkdtemp(join(tmpdir(), "tene-nextjs-"));
  const database = join(root, "orders.ndjson");
  if (!productId) {
    await rm(root, {recursive: true});
    return {screen: "empty-cart", apiStatus: 400, writes: 0};
  }
  if (options.authorized === false) {
    await rm(root, {recursive: true});
    return {screen: "forbidden", apiStatus: 403, writes: 0};
  }
  if (options.fail === true) {
    await rm(root, {recursive: true});
    return {screen: "retry-order", apiStatus: 503, writes: 0};
  }
  const order = {id: `order-${productId}`, productId};
  await import("node:fs/promises").then(({appendFile}) => appendFile(database, JSON.stringify(order) + "\n"));
  const persisted = JSON.parse((await readFile(database, "utf8")).trim());
  await rm(root, {recursive: true});
  return {screen: `confirmed:${order.id}`, apiStatus: 200, persisted, writes: 1};
}

const result = await runCheckoutJourney("sku-42");
assert.equal(result.screen, "confirmed:order-sku-42");
assert.equal(result.apiStatus, 200);
assert.deepEqual(result.persisted, {id: "order-sku-42", productId: "sku-42"});
assert.equal((await runCheckoutJourney("alternate-sku")).screen, "confirmed:order-alternate-sku");
assert.deepEqual(await runCheckoutJourney(""), {screen: "empty-cart", apiStatus: 400, writes: 0});
assert.deepEqual(await runCheckoutJourney("sku-42", {authorized: false}), {screen: "forbidden", apiStatus: 403, writes: 0});
const failed = await runCheckoutJourney("sku-retry", {fail: true});
assert.deepEqual(failed, {screen: "retry-order", apiStatus: 503, writes: 0});
const recovered = await runCheckoutJourney("sku-retry");
assert.equal(recovered.writes, 1);
assert.equal(recovered.screen, "confirmed:order-sku-retry");
console.log(JSON.stringify(result));
