import { describe, expect, it } from "vitest";

import { softThresholdBreached, usagePct } from "./budgets";

describe("budget math helpers", () => {
  it("computes usage percent", () => {
    expect(usagePct(25, 100)).toBe(25);
    expect(usagePct(10, 0)).toBe(0);
  });

  it("detects soft threshold", () => {
    expect(softThresholdBreached(79, 100, 80)).toBe(false);
    expect(softThresholdBreached(80, 100, 80)).toBe(true);
    expect(softThresholdBreached(95, 100, 80)).toBe(true);
  });
});
