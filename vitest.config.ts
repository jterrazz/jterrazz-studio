import { literate } from "@jterrazz/test/vitest";
import { defineConfig } from "vitest/config";

export default defineConfig({
  // Every `<case>.spec.yaml` under specs/ becomes a one-test module driving the
  // runner this file names — stated, never guessed.
  plugins: [literate({ specification: "./specs/cli/cli.specification.ts" })],
  test: {
    include: ["specs/**/*.test.ts"],
    testTimeout: 30_000,
    hookTimeout: 30_000,
  },
});
