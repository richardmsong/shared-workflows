# Claude Development Guidelines

This repository contains reusable GitHub Actions workflows. All workflows are designed to use GitHub App authentication instead of GITHUB_TOKEN for enhanced security and functionality.

## GitHub App Authentication

**CRITICAL**: All permissions are controlled via the Claude GitHub App. This ensures workflows can perform actions that would otherwise be blocked by GITHUB_TOKEN limitations.

### Why GitHub App Over GITHUB_TOKEN?

1. **Workflow Triggering**: PRs and commits created with GITHUB_TOKEN don't trigger other workflows. App tokens do.
2. **Enhanced Permissions**: App tokens can have more granular permissions than GITHUB_TOKEN.
3. **Security**: App tokens can be scoped per-repository with fine-grained permissions.
4. **Audit Trail**: Actions performed by the App are clearly attributed in GitHub's audit log.

## Workflow Guidelines

### 1. Never Use Permissions Blocks

**DO NOT** add `permissions:` blocks to workflow files. All permissions are managed through the GitHub App configuration.

❌ **Bad Example**:
```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    steps:
      - uses: actions/checkout@v4
```

✅ **Good Example**:
```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
```

### 2. Always Use App Token Over GITHUB_TOKEN

All GitHub API interactions and repository operations must use the Claude GitHub App token.

❌ **Bad Example**:
```yaml
- name: Checkout code
  uses: actions/checkout@v4
  with:
    token: ${{ secrets.GITHUB_TOKEN }}
```

❌ **Bad Example (Fallback Pattern)**:
```yaml
- name: Checkout code
  uses: actions/checkout@v4
  with:
    token: ${{ steps.app-token.outputs.token || github.token }}
```

✅ **Good Example**:
```yaml
- name: Generate GitHub App Token
  id: app-token
  uses: actions/create-github-app-token@v1
  with:
    app-id: ${{ secrets.CLAUDE_APP_ID }}
    private-key: ${{ secrets.CLAUDE_APP_PRIVATE_KEY }}

- name: Checkout code
  uses: actions/checkout@v4
  with:
    token: ${{ steps.app-token.outputs.token }}
```

### 3. Token Generation Pattern

Always generate the App token at the beginning of your workflow:

```yaml
jobs:
  my-job:
    runs-on: ubuntu-latest
    steps:
      - name: Generate GitHub App Token
        id: app-token
        uses: actions/create-github-app-token@v1
        with:
          app-id: ${{ secrets.CLAUDE_APP_ID }}
          private-key: ${{ secrets.CLAUDE_APP_PRIVATE_KEY }}

      # All subsequent steps should use: ${{ steps.app-token.outputs.token }}
```

### 4. Environment Variables

When using `gh` CLI or other tools that need GitHub authentication:

✅ **Good Example**:
```yaml
- name: Create PR
  env:
    GH_TOKEN: ${{ steps.app-token.outputs.token }}
  run: |
    gh pr create --title "..." --body "..."
```

## Required Changes to Existing Workflows

The following workflows need to be updated to remove permissions blocks and GITHUB_TOKEN fallbacks:

### 1. `.github/workflows/ci.yml`

**Lines 155-157**: Remove the permissions block from the `docker` job
```yaml
# REMOVE THIS:
permissions:
  contents: read
  packages: write
```

**Line 228**: Remove GITHUB_TOKEN fallback
```yaml
# CHANGE THIS:
password: ${{ secrets.REGISTRY_TOKEN || secrets.GITHUB_TOKEN }}

# TO THIS:
password: ${{ secrets.REGISTRY_TOKEN }}
```

### 2. `.github/workflows/auto-pr.yml`

**Lines 39-41**: Remove the permissions block from the `create-pr` job
```yaml
# REMOVE THIS:
permissions:
  contents: read
  pull-requests: write
```

**Lines 57 & 80**: Remove GITHUB_TOKEN/github.token fallbacks
```yaml
# CHANGE THIS:
token: ${{ steps.app-token.outputs.token || github.token }}

# TO THIS:
token: ${{ steps.app-token.outputs.token }}
```

**Line 48**: Remove `continue-on-error: true` from App token generation
```yaml
# CHANGE THIS:
- name: Generate Claude GitHub App Token
  id: app-token
  uses: actions/create-github-app-token@v1
  continue-on-error: true

# TO THIS:
- name: Generate Claude GitHub App Token
  id: app-token
  uses: actions/create-github-app-token@v1
```

### 3. `.github/workflows/tag-release.yml`

**Line 72**: Remove GITHUB_TOKEN fallback
```yaml
# CHANGE THIS:
token: ${{ secrets.RELEASE_TOKEN || secrets.GITHUB_TOKEN }}

# TO THIS:
token: ${{ secrets.RELEASE_TOKEN }}
```

### 4. `.github/workflows/claude.yml`

**Lines 44 & 51**: Remove github.token fallbacks
```yaml
# CHANGE THIS:
token: ${{ steps.app-token.outputs.token || github.token }}

# TO THIS:
token: ${{ steps.app-token.outputs.token }}
```

**Line 35**: Remove `continue-on-error: true`
```yaml
# CHANGE THIS:
- name: Generate GitHub App Token
  id: app-token
  uses: actions/create-github-app-token@v1
  continue-on-error: true

# TO THIS:
- name: Generate GitHub App Token
  id: app-token
  uses: actions/create-github-app-token@v1
```

### 5. `.github/workflows/claude-code-review.yml`

**Lines 44 & 51**: Remove github.token fallbacks
```yaml
# CHANGE THIS:
token: ${{ steps.app-token.outputs.token || github.token }}

# TO THIS:
token: ${{ steps.app-token.outputs.token }}
```

**Line 35**: Remove `continue-on-error: true`
```yaml
# CHANGE THIS:
- name: Generate GitHub App Token
  id: app-token
  uses: actions/create-github-app-token@v1
  continue-on-error: true

# TO THIS:
- name: Generate GitHub App Token
  id: app-token
  uses: actions/create-github-app-token@v1
```

## Workflow Naming Convention

Workflows that start with an underscore (`_`) are internal/test workflows specific to this repository. They exist alongside their reusable counterparts:

- `claude.yml` - Reusable workflow for other repos
- `_claude.yml` - Internal workflow that calls the reusable one

This naming convention helps distinguish between workflows meant for reuse and those used internally for testing.

## Troubleshooting

### Workflow Fails with "Resource not accessible by integration"

This error occurs when the GitHub App doesn't have the required permissions. Check:

1. App permissions in GitHub settings
2. App is installed on the repository
3. Secrets `CLAUDE_APP_ID` and `CLAUDE_APP_PRIVATE_KEY` are configured

### Workflow Fails with "Token generation failed"

If App token generation fails and there's no fallback (as intended):

1. Verify the App secrets are correctly configured
2. Check that the App is installed and not suspended
3. Review App permissions match workflow requirements

### Why No Fallback to GITHUB_TOKEN?

Fallback patterns like `|| github.token` mask configuration issues and allow workflows to run with insufficient permissions, leading to subtle failures later. By removing fallbacks, we ensure workflows fail fast and clearly when misconfigured.
