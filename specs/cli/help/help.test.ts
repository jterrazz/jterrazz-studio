import { expect, test } from "vitest";

import { cli } from "../cli.specification.js";

test("shows the full help banner when no command is given", async () => {
  // Given - no command argument
  const result = await cli.exec("");

  // Then - the whole top-level help surface matches
  expect(result.exitCode).toBe(0);
  expect(result.stdout).toMatch("help.txt");
});

test("shows the same help banner under --help", async () => {
  // Given - the explicit help flag
  const result = await cli.exec("--help");

  // Then - identical banner to the no-command form
  expect(result.exitCode).toBe(0);
  expect(result.stdout).toMatch("help.txt");
});

test("errors on an unknown command", async () => {
  // Given - a command the CLI does not know
  const result = await cli.exec("frobnicate");

  // Then - non-zero exit and a cobra error on stderr (third-party formatted, probed not snapshotted)
  expect(result.exitCode).not.toBe(0);
  expect(result.stderr).toContain("unknown command");
});
