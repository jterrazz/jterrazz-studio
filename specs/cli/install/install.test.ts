import { expect, test } from "vitest";

import { cli } from "../cli.specification.js";

test("lists the Jump Desktop tools in the catalog", async () => {
  // Given - the bare install command, which prints the tool catalog
  const result = await cli.exec("install");

  // Then - the catalog carries the expected entries. Environment-dependent: the
  // ✓/✗ install-state column and the machine header vary per host, so this is a
  // membership probe, not a full-surface golden — and a document states a stream
  // byte-exact, which is exactly what this output cannot promise.
  expect(result.exitCode).toBe(0);
  expect(result.stdout).toContain("jump-desktop-connect");
  expect(result.stdout).toContain("jump-desktop");
});
