import { test, expect } from '@playwright/test';
test.beforeEach(async ({ request }) => { await request.delete('/api/items'); });
test('success crosses UI API and persistence', async ({ page, request }) => {
  await page.goto('/'); await expect(page.getByRole('heading')).toHaveText('Create an item');
  const responsePromise=page.waitForResponse(r=>r.url().endsWith('/api/items')&&r.request().method()==='POST');
  await page.getByLabel('Item name').fill('Intent evidence'); await page.getByRole('button',{name:'Create'}).click();
  expect((await responsePromise).status()).toBe(201); await expect(page.getByRole('status')).toHaveText('Created Intent evidence');
  expect(await (await request.get('/api/state')).json()).toEqual([expect.objectContaining({name:'Intent evidence'})]);
});
test('validation prevents a network write', async ({ page, request }) => {
  await page.goto('/'); await page.getByRole('button',{name:'Create'}).click(); await expect(page.getByLabel('Item name')).toBeFocused();
  expect(await (await request.get('/api/state')).json()).toEqual([]);
});
test('downstream failure is visible and retry writes exactly once', async ({ page, request }) => {
  await page.goto('/'); await page.getByLabel('Item name').fill('Recoverable'); await page.getByLabel('Simulate one storage failure').check();
  await page.getByRole('button',{name:'Create'}).click(); await expect(page.getByRole('status')).toContainText('temporarily unavailable'); expect(await (await request.get('/api/state')).json()).toEqual([]);
  await page.getByRole('button',{name:'Retry'}).click(); await expect(page.getByRole('status')).toHaveText('Created Recoverable'); const state=await (await request.get('/api/state')).json();expect(state).toHaveLength(1);expect(state[0].name).toBe('Recoverable');
});
