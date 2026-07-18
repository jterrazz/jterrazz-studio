import { expect, test } from "vitest";

import { cli } from "../cli.specification.js";

test("documents the upgrade command under --help", async () => {
  // Given - the upgrade help flag
  const result = await cli.exec("upgrade --help");

  // Then - the whole help surface matches
  expect(result.exitCode).toBe(0);
  expect(result.stdout).toMatch("help.txt");
});
