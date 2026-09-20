import { cli, defineSpecConfig } from '@jterrazz/test/vitest';

// `cli()` is this tree's one facet: it collects `specs/cli/**/*.spec.ts` and
// wires the literate plugin onto its default runner, `specs/cli/cli.specification.ts`
// — the convention this repository already follows, so nothing is stated.
export default defineSpecConfig({
    test: {
        projects: [cli()],
    },
});
