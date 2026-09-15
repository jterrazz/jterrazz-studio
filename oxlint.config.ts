import { testing } from '@jterrazz/test/oxlint';
import { compose, defineConfig, node, type OxlintConfig } from '@jterrazz/typescript/oxlint';

/*
 * The `node` profile is the rulebook for a command-line tool, and `testing` adds
 * the spec conventions that ship with @jterrazz/test. Nothing local: this repo's
 * TypeScript is the e2e harness under `specs/`, which the shared rules already
 * describe.
 *
 * @jterrazz/test 15 declares `testing` with widened property types — its
 * override's level is `string` where oxlint takes a closed union — so the
 * fragment does not structurally satisfy `OxlintConfig`. The assertion names
 * what the fragment is until the declaration ships narrowed upstream.
 */
// oxlint-disable-next-line typescript/no-unsafe-type-assertion -- the widened upstream declaration above
export default defineConfig(compose(node, testing as OxlintConfig));
