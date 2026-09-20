import { expect, test } from 'vitest';

import { cli } from '../cli.specification.js';

// Both specs below read a surface whose VERDICTS come from the host — macOS
// version, sshd, FileVault, running services. A `<case>.spec.yaml` states a
// stream byte-exact, so these stay in code: what they prove is the presence of
// the rows the binary always emits, not the text the machine happened to fill
// them with.

const MACHINE_ROWS = ['j machine status', 'Machine', 'FileVault', 'SSH'];

const SERVICE_ROWS = [
    'Services',
    'OpenClaw runtime',
    'OpenClaw config',
    'Hermes runtime',
    'Hermes config',
    'OrbStack',
];

test('reports the machine-level state checks', async () => {
    // Given - the machine status command on the default (client) host
    const result = await cli.exec('machine status');

    // Then - the Machine section and its stable row labels are present
    expect(result.exitCode).toBe(0);
    for (const row of MACHINE_ROWS) {
        expect(result.stdout).toContain(row);
    }
});

test('reports the service state checks on a server machine', async () => {
    // Given - a registry declaring this host a server, which gates the Services section
    const result = await cli
        .fixture('server-registry/')
        .env({ HOME: '$WORKDIR' })
        .exec('machine status');

    // Then - every server-only service row is present
    expect(result.exitCode).toBe(0);
    for (const row of SERVICE_ROWS) {
        expect(result.stdout).toContain(row);
    }
});
