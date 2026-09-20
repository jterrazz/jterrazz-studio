import { testing } from '@jterrazz/test/oxlint';
import { compose, defineConfig, node } from '@jterrazz/typescript/oxlint';

/*
 * The `node` profile is the rulebook for a command-line tool, and `testing` adds
 * the spec conventions that ship with @jterrazz/test. Nothing local: this repo's
 * TypeScript is the e2e harness under `specs/`, which the shared rules already
 * describe.
 */
export default defineConfig(compose(node, testing));
