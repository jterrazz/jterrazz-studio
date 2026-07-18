import { expect, test } from "vitest";

import { cli } from "../cli.specification.js";

test("lists the machine subcommands under --help", async () => {
  // Given - the machine command group help flag
  const result = await cli.exec("machine --help");

  // Then - the whole subcommand listing matches
  expect(result.exitCode).toBe(0);
  expect(result.stdout).toMatch("help.txt");
});

test("documents preboot SSH under machine unlock --help", async () => {
  // Given - the machine unlock help flag
  const result = await cli.exec("machine unlock --help");

  // Then - the whole help surface matches
  expect(result.exitCode).toBe(0);
  expect(result.stdout).toMatch("unlock.txt");
});

test("reports the machine-level state checks", async () => {
  // Given - the machine status command on the default (client) host
  const result = await cli.exec("machine status");

  // Then - the Machine section and its stable row labels are present.
  // Environment-dependent: the ✅/❌ verdicts and their detail text vary per host
  // (macOS version, sshd, FileVault), so this probes the stable row labels the
  // binary always emits, not the whole surface.
  expect(result.exitCode).toBe(0);
  expect(result.stdout).toContain("j machine status");
  expect(result.stdout).toContain("Machine");
  expect(result.stdout).toContain("FileVault");
  expect(result.stdout).toContain("SSH");
});

test("reports the service state checks on a server machine", async () => {
  // Given - a registry that declares this host a server, so the Services section
  // (gated on the server role) is emitted deterministically regardless of host
  const result = await cli
    .fixture("server-registry/")
    .env({ HOME: "$WORKDIR" })
    .exec("machine status");

  // Then - the server-only service rows are present. Same rationale as above:
  // the verdicts and details vary per host, so this probes the stable product
  // surface — the row labels the binary emits in a deterministic order.
  expect(result.exitCode).toBe(0);
  expect(result.stdout).toContain("Services");
  expect(result.stdout).toContain("OpenClaw runtime");
  expect(result.stdout).toContain("OpenClaw config");
  expect(result.stdout).toContain("Hermes runtime");
  expect(result.stdout).toContain("Hermes config");
  expect(result.stdout).toContain("OrbStack");
});
