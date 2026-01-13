# Shared GitHub Workflows

Reusable GitHub Actions workflows for Go/kubebuilder projects that publish Docker images to GHCR, enabling downstream repositories to maintain consistent CI/CD pipelines.

## Available Workflows

| Workflow | Description |
|----------|-------------|
| [`ci.yml`](#ciyml) | Main CI workflow with test, lint, build, and Docker stages |
| [`tag-release.yml`](#tag-releaseyml) | Release automation with semver tagging and version branches |
| [`auto-pr.yml`](#auto-pryml) | Automatic PR creation for `claude/**` branches |
| [`claude.yml`](#claudeyml) | Claude AI interaction on issues/PRs/comments |
| [`claude-code-review.yml`](#claude-code-reviewyml) | Automated code review via Claude |

## Workflow Naming Convention

### Underscore Prefix for Internal Workflows

Workflows prefixed with an underscore (`_`) are **internal workflows** that run within this repository itself. These workflows are used for testing and managing the shared-workflows repository, not for use by downstream repositories.

**Convention:**
- **Public/Reusable workflows** (no underscore): Called by downstream repositories via `uses: richardmsong/shared-workflows/.github/workflows/<name>.yml@v1`
- **Internal workflows** (underscore prefix): Run only within this repository for testing and internal automation

**Examples:**

| Public Workflow | Internal Workflow | Purpose |
|----------------|-------------------|---------|
| `ci.yml` | `_test-ci.yml` | Tests the reusable CI workflow using test fixtures |
| `tag-release.yml` | `_test-release.yml` | Tests the release workflow in dry-run mode |
| `claude.yml` | `_claude.yml` | Triggers Claude AI for this repository's issues/PRs |
| `claude-code-review.yml` | `_claude-code-review.yml` | Enables Claude code review for this repository |

This naming convention makes it immediately clear which workflows are intended for external use and which are for internal repository management.

## Quick Start

### CI Workflow

Create `.github/workflows/ci.yml` in your repository:

```yaml
name: CI

on:
  push:
    branches: [main, master]
    tags: ['v*.*.*']
  pull_request:
    branches: [main, master]

jobs:
  ci:
    uses: richardmsong/shared-workflows/.github/workflows/ci.yml@v1
    secrets: inherit
```

### Release Workflow

Create `.github/workflows/release.yml` in your repository:

```yaml
name: Release

on:
  release:
    types: [published]
  workflow_dispatch:
    inputs:
      version:
        description: 'Version to release (e.g., 1.0.0)'
        required: true
      dry_run:
        description: 'Dry run mode'
        type: boolean
        default: true

jobs:
  release:
    uses: richardmsong/shared-workflows/.github/workflows/tag-release.yml@v1
    with:
      version: ${{ inputs.version || '' }}
      dry-run: ${{ inputs.dry_run || false }}
      image-name: ghcr.io/${{ github.repository }}
    secrets: inherit
```

### Auto PR Workflow

Create `.github/workflows/auto-pr.yml` in your repository:

```yaml
name: Auto Create PR

on:
  push:
    branches:
      - 'claude/**'

jobs:
  auto-pr:
    uses: richardmsong/shared-workflows/.github/workflows/auto-pr.yml@v1
    secrets: inherit
```

### Claude Interaction Workflow

Create `.github/workflows/claude.yml` in your repository:

```yaml
name: Claude Code

on:
  issue_comment:
    types: [created]
  pull_request_review_comment:
    types: [created]
  issues:
    types: [opened, assigned]
  pull_request_review:
    types: [submitted]

jobs:
  claude:
    if: |
      (github.event_name == 'issue_comment' && contains(github.event.comment.body, '@claude')) ||
      (github.event_name == 'pull_request_review_comment' && contains(github.event.comment.body, '@claude')) ||
      (github.event_name == 'pull_request_review' && contains(github.event.review.body, '@claude')) ||
      (github.event_name == 'issues' && (contains(github.event.issue.body, '@claude') || contains(github.event.issue.title, '@claude')))
    uses: richardmsong/shared-workflows/.github/workflows/claude.yml@v1
    secrets: inherit
```

### Claude Code Review Workflow

Create `.github/workflows/claude-code-review.yml` in your repository:

```yaml
name: Claude Code Review

on:
  pull_request:
    types: [opened, synchronize]

jobs:
  review:
    uses: richardmsong/shared-workflows/.github/workflows/claude-code-review.yml@v1
    secrets: inherit
```

## Workflow Reference

### ci.yml

Main CI workflow providing test, lint, build, and Docker image publishing stages.

#### Inputs

| Input | Type | Default | Description |
|-------|------|---------|-------------|
| `working-directory` | string | `.` | Working directory for the project |
| `go-version` | string | `1.22` | Go version (ignored if go-version-file is set) |
| `go-version-file` | string | `go.mod` | Path to go.mod for Go version detection |
| `registry` | string | `ghcr.io` | Container registry URL |
| `image-name` | string | `''` | Docker image name (defaults to repository name) |
| `push-image` | boolean | `true` | Whether to push the Docker image |
| `run-tests` | boolean | `true` | Whether to run tests |
| `run-lint` | boolean | `true` | Whether to run linting |
| `lint-version` | string | `v2.7.2` | golangci-lint version |
| `dockerfile` | string | `Dockerfile` | Path to Dockerfile |
| `docker-build-target` | string | `docker-buildx-ci` | Makefile target for Docker build |

#### Outputs

| Output | Description |
|--------|-------------|
| `image-tags` | Docker image tags that were built |

#### Secrets

| Secret | Required | Description |
|--------|----------|-------------|
| `REGISTRY_TOKEN` | No | Token for container registry (defaults to GITHUB_TOKEN) |

#### Docker Tags Generated

The workflow generates intelligent Docker tags:

- `sha-<short-sha>` - Always added for ephemeral builds
- `<branch-name>` - Added for branch pushes
- `pr-<number>` - Added for pull requests
- `<version>`, `<major>.<minor>`, `<major>` - Added for semver tags (v*.*.*)
- `latest` - Added for default branch and semver tags

### tag-release.yml

Release automation workflow with semver tagging and version branch management.

#### Inputs

| Input | Type | Default | Description |
|-------|------|---------|-------------|
| `version` | string | `''` | Version to release (e.g., 1.0.0) |
| `dry-run` | boolean | `false` | Perform dry run without making changes |
| `image-name` | string | `''` | Docker image name for kustomization update |
| `kustomization-path` | string | `config/manager` | Path to kustomization.yaml directory |
| `update-kustomization` | boolean | `true` | Whether to update kustomization.yaml |
| `create-major-branch` | boolean | `true` | Create/update major version branch (e.g., v1) |
| `create-minor-branch` | boolean | `true` | Create/update minor version branch (e.g., v1.2) |

#### Outputs

| Output | Description |
|--------|-------------|
| `version` | The released version |
| `tag` | The created tag |

#### Secrets

| Secret | Required | Description |
|--------|----------|-------------|
| `RELEASE_TOKEN` | No | Token for pushing tags/branches (defaults to GITHUB_TOKEN) |

### auto-pr.yml

Automatically creates pull requests for `claude/**` branches.

#### Inputs

| Input | Type | Default | Description |
|-------|------|---------|-------------|
| `branch-pattern` | string | `^claude/` | Branch pattern regex to trigger PR creation |
| `base-branch` | string | `''` | Base branch for PR (defaults to repository default) |

#### Secrets

| Secret | Required | Description |
|--------|----------|-------------|
| `CLAUDE_APP_ID` | No | Claude GitHub App ID |
| `CLAUDE_APP_PRIVATE_KEY` | No | Claude GitHub App private key |

#### Prerequisites

1. Enable "Allow GitHub Actions to create and approve pull requests" in repository settings:
   - Settings → Actions → General → Workflow permissions
2. (Optional) Configure Claude GitHub App secrets for PR events to trigger other workflows

### claude.yml

Claude AI interaction workflow for issues, PRs, and comments.

#### Inputs

| Input | Type | Default | Description |
|-------|------|---------|-------------|
| `allowed-tools` | string | See below | Allowed tools for Claude Code |
| `additional-permissions` | string | `actions: read` | Additional GitHub API permissions |

Default allowed tools: `Bash(gh:*),Bash(go:*),Bash(make:*),Bash(git:*),Bash(npm:*),Bash(cargo:*),mcp__codesign__sign_file,mcp__github__*,WebFetch,WebSearch`

> **Note:** The `mcp__github__*` wildcard grants blanket access to all GitHub MCP tools. This is intentional because the Claude workflow is designed for read-write operations (creating issues, PRs, commits, etc.). Real security is enforced by your GitHub App permissions, not this allowlist.

#### Secrets

| Secret | Required | Description |
|--------|----------|-------------|
| `CLAUDE_APP_ID` | No | Claude GitHub App ID |
| `CLAUDE_APP_PRIVATE_KEY` | No | Claude GitHub App private key |
| `CLAUDE_CODE_OAUTH_TOKEN` | **Yes** | Claude Code OAuth token |

### claude-code-review.yml

Automated code review workflow using Claude.

#### Inputs

| Input | Type | Default | Description |
|-------|------|---------|-------------|
| `review-prompt` | string | `''` | Custom review prompt (optional) |
| `allowed-tools` | string | See below | Allowed tools for Claude Code review |

Default allowed tools:
```
mcp__github__create_pending_pull_request_review,
mcp__github__add_comment_to_pending_review,
mcp__github__submit_pending_pull_request_review,
mcp__github__get_pull_request_diff,
mcp__github__get_pull_request_files,
mcp__github__get_pull_request_review_comments,
mcp__github__get_pull_request,
mcp__github__list_pull_requests,
mcp__github__search_code,
mcp__github__search_issues,
Bash(gh pr checks:*),
Bash(gh pr view:*),
Bash(gh:*)
```

> **Note:** Unlike the main `claude.yml` workflow, this review workflow uses granular MCP permissions because it's designed for read-only code review operations. It should NOT make changes to code, branches, or repository state.

#### Secrets

| Secret | Required | Description |
|--------|----------|-------------|
| `CLAUDE_APP_ID` | No | Claude GitHub App ID |
| `CLAUDE_APP_PRIVATE_KEY` | No | Claude GitHub App private key |
| `CLAUDE_CODE_OAUTH_TOKEN` | **Yes** | Claude Code OAuth token |

## Required Repository Setup

### Makefile Targets

The CI workflow expects these Makefile targets:

```makefile
test:       # Run tests, output coverage to cover.out
build:      # Build the Go binary
docker-buildx-ci:  # Build multi-arch Docker image
            # Should use DOCKER_TAGS and DOCKER_PUSH env vars
```

See `test-fixtures/kubebuilder-minimal/Makefile` for a reference implementation.

### Secrets Configuration

Configure these secrets in your repository:

| Secret | Workflows | Description |
|--------|-----------|-------------|
| `CLAUDE_CODE_OAUTH_TOKEN` | claude, claude-code-review | OAuth token for Claude Code |
| `CLAUDE_APP_ID` | auto-pr, claude, claude-code-review | Claude GitHub App ID (optional) |
| `CLAUDE_APP_PRIVATE_KEY` | auto-pr, claude, claude-code-review | Claude GitHub App private key (optional) |

### Permissions

Ensure your workflows have appropriate permissions:

```yaml
permissions:
  contents: read
  packages: write      # For Docker push
  pull-requests: write # For auto-pr
```

## SDLC Practices

### Versioning

This repository follows semantic versioning:

- Use `@v1` (major) for production stability
- Use `@v1.2` (minor) for specific feature sets
- Use `@v1.2.3` (patch) for exact version pinning

### Testing

Internal test workflows validate changes:

- `_test-ci.yml` - Tests CI workflow against fixture project
- `_test-release.yml` - Tests release workflow in dry-run mode

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes and test locally
4. Submit a pull request

## License

MIT License - See [LICENSE](LICENSE) for details.
