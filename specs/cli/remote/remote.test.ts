import { expect, test } from "vitest";

import { cli } from "../cli.specification.js";

test("lists the remote access subcommands under --help", async () => {
  // Given - the remote command group help flag
  const result = await cli.exec("remote --help");

  // Then - the whole subcommand listing matches
  expect(result.exitCode).toBe(0);
  expect(result.stdout).toMatch("help.txt");
});
