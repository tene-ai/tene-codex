import { defineConfig } from '@playwright/test';
export default defineConfig({
  testDir: './tests/e2e', workers: 1, fullyParallel: false, reporter: [['list'], ['json', { outputFile: 'test-results/playwright.json' }]],
  use: { baseURL: 'http://127.0.0.1:4173', trace: 'retain-on-failure', screenshot: 'only-on-failure' },
  webServer: { command: 'go run ./testdata/reference-web', url: 'http://127.0.0.1:4173', reuseExistingServer: false, env: { TENE_REFERENCE_STATE: '.tmp/reference-state.json' } }
});
