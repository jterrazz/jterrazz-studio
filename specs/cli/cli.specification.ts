import { execSync } from "node:child_process";
import { existsSync, readdirSync, statSync } from "node:fs";
import { resolve } from "node:path";

import { specification } from "@jterrazz/test";
import { afterAll } from "vitest";

const REPO_ROOT = resolve(import.meta.dirname, "../..");
const J_BIN = resolve(REPO_ROOT, ".artifacts/go/j-test");
const GO_SRC = resolve(REPO_ROOT, "src");

// A stale binary silently tests old behaviour. Rebuild when it is missing, when
// a rebuild is forced, or when any `src/**/*.go` is newer than the binary — an
// mtime check so editing the CLI and re-running the suite exercises the change.
function binaryIsStale(): boolean {
  if (!existsSync(J_BIN)) {
    return true;
  }
  const binModifiedMs = statSync(J_BIN).mtimeMs;
  for (const entry of readdirSync(GO_SRC, { recursive: true, withFileTypes: true })) {
    if (!entry.isFile() || !entry.name.endsWith(".go")) {
      continue;
    }
    if (statSync(resolve(entry.parentPath, entry.name)).mtimeMs > binModifiedMs) {
      return true;
    }
  }
  return false;
}

// The ONE real runner: every spec drives the built `j` binary. Build it once,
// before any spec runs, so the whole suite shares a single compilation.
if (process.env.J_FORCE_REBUILD === "1" || binaryIsStale()) {
  execSync(`go build -o ${J_BIN} ./src/cmd/j`, {
    cwd: REPO_ROOT,
    stdio: "inherit",
  });
}

export const { cleanup, cli } = await specification.cli(J_BIN);

afterAll(cleanup);
