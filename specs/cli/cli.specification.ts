import { specification } from '@jterrazz/test';
import { execSync } from 'node:child_process';
import { existsSync, readdirSync, statSync } from 'node:fs';
import { resolve } from 'node:path';
import { afterAll } from 'vitest';

const REPO_ROOT = resolve(import.meta.dirname, '../..');
const WORKSTATION_BIN = resolve(REPO_ROOT, '.artifacts/go/workstation-test');
const GO_SRC = resolve(REPO_ROOT, 'src');

// hazard: a stale binary silently tests old behaviour, so the suite rebuilds when
// any `src/**/*.go` is newer than the binary, so an edit to the CLI is what a rerun exercises.
function binaryIsStale(): boolean {
    if (!existsSync(WORKSTATION_BIN)) {
        return true;
    }
    const binModifiedMs = statSync(WORKSTATION_BIN).mtimeMs;
    for (const entry of readdirSync(GO_SRC, { recursive: true, withFileTypes: true })) {
        if (!entry.isFile() || !entry.name.endsWith('.go')) {
            continue;
        }
        if (statSync(resolve(entry.parentPath, entry.name)).mtimeMs > binModifiedMs) {
            return true;
        }
    }
    return false;
}

if (process.env.WORKSTATION_FORCE_REBUILD === '1' || binaryIsStale()) {
    execSync(`go build -o ${WORKSTATION_BIN} ./src/cmd/workstation`, {
        cwd: REPO_ROOT,
        stdio: 'inherit',
    });
}

// invariant: `workstation` reads its machine registry from `~/.studio/`, so every run
// gets a HOME of its own and never the developer's.
export const { cleanup, cli } = await specification.cli(WORKSTATION_BIN, {
    defaults: { HOME: '$WORKDIR' },
});

afterAll(cleanup);
