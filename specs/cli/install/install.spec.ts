import { expect, test } from 'vitest';

import { cli } from '../cli.specification.js';

// This spec reads a surface whose ✓/✗ install-state column and machine header
// vary per host. A `<case>.spec.yaml` states a stream byte-exact, so it stays in
// code: what it proves is the presence of the catalog's entries, not the text the
// machine happened to fill them with.

test('lists the Jump Desktop tools in the catalog', async () => {
    // Given - the bare install command, which prints the tool catalog
    const result = await cli.exec('install');

    // Then - the catalog carries both Jump Desktop entries
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toContain('jump-desktop-connect');
    expect(result.stdout).toContain('jump-desktop');
});
