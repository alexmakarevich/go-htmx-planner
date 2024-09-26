import { test, expect } from "@playwright/test";
import { BASE_URL } from "./utils";
import { v4 } from "uuid";

const testId = v4();
const testEventTitle = "Test-Event-" + testId;
const testEventTitleUpdated = testEventTitle + "- Edited";

// TODO: get elements in the most semantic/stable way possible

test("cycles CRUD of event", async ({ page }) => {
  await page.goto(BASE_URL);

  // go to events page
  await page.getByText(/events/).click();

  // go to events page
  await page.getByText(/create/).click();

  await page.locator("input[name='title']").fill(testEventTitle);

  // try to create
  await page.getByText(/create/).click();

  // cannot create event without date
  await expect(page.getByText(/validation.*failed/)).toBeVisible();

  await page.locator("input[name='date-time']").fill("2024-09-26T21:14");

  // create
  await page.getByText(/create/).click();

  // expect redirect to events page
  await expect(page.getByText(/Events/)).toBeVisible();

  const event = page.getByRole("listitem").filter({ hasText: testEventTitle });

  // edit event
  await event.getByText(/edit/).click();

  await page.locator("input[name='title']").fill(testEventTitleUpdated);
  await page.locator("input[name='date-time']").fill("2024-10-10T10:00");
  await page.getByText(/update/).click();

  // expect redirect to events page
  await expect(page.getByText(/Events/)).toBeVisible();
  // get the event
  const eventAfterUpdate = page
    .getByRole("listitem")
    .filter({ hasText: testEventTitleUpdated });

  await expect(eventAfterUpdate).toContainText(/10 Oct 2024/);

  // delete event
  await event.getByText(/delete/).click();

  await expect(
    page.getByRole("listitem").filter({ hasText: testEventTitleUpdated })
  ).not.toBeVisible();
});
