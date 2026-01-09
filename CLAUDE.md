# Claude Instructions for shared-workflows Repository

## Permissions and Capabilities

### Workflow File Editing
**You CAN edit workflow files in the `.github/workflows/` directory.**

This repository uses a custom GitHub App configuration that provides Claude with permissions to modify workflow files, unlike the default Anthropic Claude Code GitHub App. When working with CI/CD workflows, you can:

- Edit existing workflow files in `.github/workflows/`
- Create new workflow files
- Update workflow configurations, job definitions, and steps
- Modify reusable workflow definitions

Do not assume you lack permissions to edit workflow files - you have full access to make necessary changes.

## Repository Structure

This repository contains reusable GitHub Actions workflows:

- `.github/workflows/ci.yml` - Main reusable CI workflow for Go projects
- `.github/workflows/_test-ci.yml` - Test workflow that validates the CI workflow
- `test-fixtures/` - Test fixtures used to validate workflows

## GitHub App Authentication Requirements

**CRITICAL: All workflows must use GitHub App authentication exclusively.**

All permissions are controlled via the GitHub App. To ensure workflows function correctly:

### DO NOT use permissions blocks
- **Never** add `permissions:` blocks to job definitions
- All required permissions are granted through the GitHub App
- Permissions blocks cause workflows to fail with "Requested permissions" errors

### DO NOT use GITHUB_TOKEN
- **Never** use `secrets.GITHUB_TOKEN` as a fallback
- **Never** use `github.token` as a fallback
- All GitHub API interactions must use tokens from the GitHub App
- Use `${{ secrets.CLAUDE_APP_ID }}` and `${{ secrets.CLAUDE_APP_PRIVATE_KEY }}` with `actions/create-github-app-token@v1`

### DO NOT use continue-on-error for App token generation
- App token generation steps must not have `continue-on-error: true`
- If the App token generation fails, the workflow should fail
- This ensures we never silently fall back to GITHUB_TOKEN

### Correct Pattern for GitHub App Authentication

```yaml
steps:
  - name: Generate GitHub App Token
    id: app-token
    uses: actions/create-github-app-token@v1
    with:
      app-id: ${{ secrets.CLAUDE_APP_ID }}
      private-key: ${{ secrets.CLAUDE_APP_PRIVATE_KEY }}

  - name: Checkout repository
    uses: actions/checkout@v4
    with:
      token: ${{ steps.app-token.outputs.token }}

  - name: Use gh CLI or GitHub API
    env:
      GH_TOKEN: ${{ steps.app-token.outputs.token }}
    run: gh pr create ...
```

### Incorrect Patterns (DO NOT USE)

```yaml
# ❌ WRONG - permissions block
jobs:
  my-job:
    permissions:
      contents: read
      packages: write

# ❌ WRONG - GITHUB_TOKEN fallback
token: ${{ steps.app-token.outputs.token || secrets.GITHUB_TOKEN }}

# ❌ WRONG - github.token fallback
token: ${{ steps.app-token.outputs.token || github.token }}

# ❌ WRONG - continue-on-error on token generation
- name: Generate GitHub App Token
  continue-on-error: true
  uses: actions/create-github-app-token@v1
```

## Important Notes

- When updating workflows, ensure changes are tested via the `_test-ci.yml` workflow
- Follow GitHub Actions best practices for reusable workflows
- All workflows must use GitHub App authentication exclusively - no GITHUB_TOKEN fallbacks
- Never add permissions blocks to jobs - permissions are controlled by the GitHub App
