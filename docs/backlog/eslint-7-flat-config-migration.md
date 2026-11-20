---
worth: later
where: hugo/.eslintrc.js
added: 2026-11-20
---
# eslint 7 and typescript-eslint 4 need a flat-config migration to move forward

Dependabot PR #516 bumps eslint 7 -> 10 and @typescript-eslint/eslint-plugin 4 -> 8 (plus js-yaml through
the same tree). ESLint 10 removes legacy `.eslintrc` support, and size, frontend-lint and build-and-compare
all fail on that PR while they pass on master. The upgrade needs `hugo/.eslintrc.js` rewritten as
`eslint.config.js`, the typescript-eslint and eslint-config-prettier packages moved to versions that
support flat config, and `.eslintignore` folded into the config's `ignores`. The size and site-build
failures run through Mix, so check them separately once lint passes.
