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

## Authentication in Workflows

This repository's workflows use different authentication methods depending on the operation:

### For Claude's Operations (Editing Files, Creating PRs, etc.)
**Use GitHub App authentication exclusively.**

When Claude needs to interact with the repository (creating PRs, pushing commits to branches, etc.), use the GitHub App token:

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

**DO NOT use continue-on-error for App token generation:**
- App token generation steps must not have `continue-on-error: true`
- If the App token generation fails, the workflow should fail
- This ensures we never silently fall back to GITHUB_TOKEN

### For Repository Operations (Pushing Packages, etc.)
**Use the workflow's own GITHUB_TOKEN.**

For operations where the repository itself should be the actor (like pushing Docker images to GHCR), use `secrets.GITHUB_TOKEN`:

```yaml
jobs:
  docker:
    permissions:
      contents: read
      packages: write
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Log in to registry
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
```

**Rationale:** Package operations should appear as coming from the repository itself, not from Claude's GitHub App. This maintains clear ownership and permissions boundaries.

### Using permissions blocks
- You **may** add `permissions:` blocks to job definitions when needed
- Permissions blocks are required for certain operations like pushing packages
- The permissions you request apply to the `GITHUB_TOKEN` for that job

### Incorrect Patterns (DO NOT USE)

```yaml
# ❌ WRONG - GITHUB_TOKEN fallback for App operations
token: ${{ steps.app-token.outputs.token || secrets.GITHUB_TOKEN }}

# ❌ WRONG - continue-on-error on token generation
- name: Generate GitHub App Token
  continue-on-error: true
  uses: actions/create-github-app-token@v1

# ❌ WRONG - Using App token for package operations
# (packages should be pushed by the repository, not Claude's App)
password: ${{ steps.app-token.outputs.token }}
```

## Important Notes

- When updating workflows, ensure changes are tested via the `_test-ci.yml` workflow
- Follow GitHub Actions best practices for reusable workflows
- Use GitHub App authentication for Claude's operations (PRs, commits, etc.)
- Use `secrets.GITHUB_TOKEN` for repository operations (pushing packages, etc.)
- Add permissions blocks when needed for operations like pushing packages
