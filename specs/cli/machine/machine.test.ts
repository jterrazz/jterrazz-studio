import { expect, test } from "vitest";

import { cli } from "../cli.specification.js";

// Both specs below read a surface whose VERDICTS come from the host — macOS
// version, sshd, FileVault, running services. A `<case>.spec.yaml` states a
// stream byte-exact, so these stay in code: what they prove is the presence of
// the rows the binary always emits, not the text the machine happened to fill
// them with.

test("reports the machine-level state checks", async () => {
  // Given - the machine status command on the default (client) host
  const result = await cli.exec("machine status");

  // Then - the Machine section and its stable row labels are present
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

  // Then - the server-only service rows are present, in the deterministic order
  // the binary emits them
  expect(result.exitCode).toBe(0);
  expect(result.stdout).toContain("Services");
  expect(result.stdout).toContain("OpenClaw runtime");
  expect(result.stdout).toContain("OpenClaw config");
  expect(result.stdout).toContain("Hermes runtime");
  expect(result.stdout).toContain("Hermes config");
  expect(result.stdout).toContain("OrbStack");
});
