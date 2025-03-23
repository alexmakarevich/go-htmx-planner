import { test, expect } from "@playwright/test";
import { BASE_URL } from "./utils";

test("basic", async ({ page }) => {
  await page.goto(BASE_URL);

  // Expect a title "to contain" a substring.
  await expect(page).toHaveTitle(/Go-Htmx-Planner/);
});
