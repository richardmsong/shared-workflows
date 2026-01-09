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

### Understanding github.token
- The `github.token` context variable represents **repository delegated permissions** from the GitHub App
- When a GitHub App is installed on a repository, `github.token` contains the App's delegated token, not the default GITHUB_TOKEN
- This means `github.token` is **safe to use** and represents the GitHub App's permissions
- See [GitHub's documentation on repository delegated permissions](https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/about-authentication-with-a-github-app#repository-delegated-permissions)

### Using permissions blocks
- You **may** add `permissions:` blocks to job definitions when needed
- Permissions blocks are required for certain operations like pushing packages
- The permissions you request must be within the scope granted to the GitHub App

### DO NOT use GITHUB_TOKEN secret
- **Never** use `secrets.GITHUB_TOKEN` as it bypasses the GitHub App
- Always use `github.token` (delegated App permissions) or generate tokens via `actions/create-github-app-token@v1`

### DO NOT use continue-on-error for App token generation
- App token generation steps must not have `continue-on-error: true`
- If the App token generation fails, the workflow should fail
- This ensures we never silently fall back to GITHUB_TOKEN

### Correct Patterns for GitHub App Authentication

#### Using github.token (delegated App permissions)
```yaml
jobs:
  my-job:
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
          password: ${{ github.token }}
```

#### Generating explicit App token
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
# ❌ WRONG - GITHUB_TOKEN secret
token: ${{ secrets.GITHUB_TOKEN }}

# ❌ WRONG - GITHUB_TOKEN fallback
token: ${{ steps.app-token.outputs.token || secrets.GITHUB_TOKEN }}

# ❌ WRONG - continue-on-error on token generation
- name: Generate GitHub App Token
  continue-on-error: true
  uses: actions/create-github-app-token@v1
```

## Important Notes

- When updating workflows, ensure changes are tested via the `_test-ci.yml` workflow
- Follow GitHub Actions best practices for reusable workflows
- All workflows must use GitHub App authentication exclusively - never use `secrets.GITHUB_TOKEN`
- Use `github.token` for delegated App permissions (safe and recommended)
- Add permissions blocks when needed for operations like pushing packages
