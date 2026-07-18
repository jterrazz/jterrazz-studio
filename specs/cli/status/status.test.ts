import { expect, test } from "vitest";

import { cli } from "../cli.specification.js";

test("documents the status command under --help", async () => {
  // Given - the status help flag
  const result = await cli.exec("status --help");

  // Then - the whole help surface matches
  expect(result.exitCode).toBe(0);
  expect(result.stdout).toMatch("help.txt");
});
