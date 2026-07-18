// Conventions enforcement, standalone (docs/10 — "Standalone — without
// @jterrazz/typescript"). This repo is a Go CLI with a thin TS test suite and
// keeps its OWN formatting (2-space indent, double quotes) — so it adopts only
// the `testing` fragment (the whole jterrazz/* catalogue + the A4 overrides)
// and does NOT pull in @jterrazz/typescript's formatting/base preset.
import { testing } from "@jterrazz/test/oxlint";

export default testing;
