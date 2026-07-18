import { expect, test } from "vitest";

import { cli } from "../cli.specification.js";

test("lists the git workflow subcommands under --help", async () => {
  // Given - the git command group help flag
  const result = await cli.exec("run git --help");

  // Then - the whole subcommand listing matches
  expect(result.exitCode).toBe(0);
  expect(result.stdout).toMatch("git.txt");
});

test("lists the docker workflow subcommands under --help", async () => {
  // Given - the docker command group help flag
  const result = await cli.exec("run docker --help");

  // Then - the whole subcommand listing matches
  expect(result.exitCode).toBe(0);
  expect(result.stdout).toMatch("docker.txt");
});
