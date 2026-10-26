import { test } from '@jterrazz/test/config';
import { defineConfig } from '@jterrazz/typescript/config';

export default defineConfig({
    kind: 'cli',
    platform: 'node',
    presets: [test()],
    typescript: {
        knip: {
            entry: {
                'specs/cli/cli.specification.ts':
                    'the runner every `<case>.spec.yaml` document is bound to by the literate plugin; no TypeScript spec imports it',
            },
            ignoreBinaries: {
                go: 'the CLI specification builds the Go binary it drives; the toolchain is a machine requirement, never an npm dependency',
            },
        },
        roots: {
            dotfiles:
                'the application configs `workstation config` installs onto a machine, versioned with the CLI that writes them',
            scripts:
                'the curl-piped installer a fresh machine runs before `workstation` exists on it',
        },
    },
});
