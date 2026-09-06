import { defineSpecConfig } from "@jterrazz/test/vitest";

export default defineSpecConfig({
  // Every `<case>.spec.yaml` under specs/ becomes a one-test module driving the
  // runner this file names — stated, never guessed.
  literate: { specification: "./specs/cli/cli.specification.ts" },
  test: {
    include: ["specs/**/*.test.ts"],
  },
});
