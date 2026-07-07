import { dirname, resolve } from "node:path";
import { fileURLToPath, pathToFileURL } from "node:url";
import { createRequire } from "node:module";

const require = createRequire(import.meta.url);
const { chromium } = require("playwright");

const dir = dirname(fileURLToPath(import.meta.url));
const htmlPath = resolve(dir, "cake-go-masker.html");

const browser = await chromium.launch();
const page = await browser.newPage({
  viewport: { width: 1400, height: 1500 },
  deviceScaleFactor: 1,
});

await page.goto(pathToFileURL(htmlPath).href);
await page.evaluate(() => document.fonts.ready);

const targets = await page.locator("[data-file]").evaluateAll((nodes) =>
  nodes.map((node) => node.getAttribute("data-file"))
);

for (const file of targets) {
  const target = page.locator(`[data-file="${file}"]`);
  await target.screenshot({ path: resolve(dir, file) });
}

await browser.close();
