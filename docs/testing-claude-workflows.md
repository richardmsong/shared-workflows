# Testing Claude Workflows

This document describes how to test Claude and Claude Code Review functionality in this repository.

## Overview

The repository provides reusable workflows for Claude AI integration. These workflows must be properly configured to ensure Claude can:

1. **Pick up @claude mentions** on issues, PRs, and comments
2. **Trigger workflows** correctly when notified
3. **Use allowed tools** (go, npm, bash, GitHub MCP, etc.)

## Automated Tests

### Workflow Structure Validation

The `_test-claude-workflows.yml` workflow validates the structure and configuration of Claude workflows:

**What it checks:**
- ✅ Workflows use `workflow_call` trigger (not direct event triggers)
- ✅ `allowed-tools` input exists with sensible defaults
- ✅ Inputs are properly passed to Claude Code action
- ✅ Required secrets are documented
- ✅ Internal wrapper workflows correctly call reusable workflows
- ✅ Tool patterns follow expected format

**When it runs:**
- Automatically on PRs that modify Claude workflow files
- Can be triggered manually via workflow_dispatch

**Why this matters:**

The test catches the exact issue described in [#5](https://github.com/richardmsong/shared-workflows/issues/5):
- Using `pull_request` events instead of `workflow_call` leaves inputs blank
- This causes `allowedTools` to be empty, breaking Claude's functionality
- The test fails if workflows use direct event triggers instead of `workflow_call`

## Manual Testing

### Testing Claude on Issues

1. Create a test issue in your repository
2. Include `@claude` in the issue body or title
3. Verify that the Claude workflow triggers
4. Check the workflow run to ensure:
   - Checkout succeeds
   - Claude Code action runs
   - Claude has access to required tools
   - Claude can update its comment

### Testing Claude Code Review

1. Open a test pull request
2. Verify that the Claude Code Review workflow triggers
3. Check that Claude:
   - Can read the PR diff
   - Can use `gh` commands to interact with GitHub
   - Can post review comments
   - Has access to configured tools

### Testing Tool Access

To verify Claude has proper tool access, ask it to run specific commands:

```
@claude Can you run `go version` and tell me what version is available?
```

Expected behavior:
- Claude should be able to execute the command if `Bash(go:*)` is in allowed-tools
- If the tool isn't allowed, Claude should report it cannot access that tool

## Common Issues

### Issue: Claude cannot update its comment

**Symptoms:**
- Claude workflow runs but comment doesn't update
- Permission denied errors in workflow logs

**Root cause:**
- Missing or incorrect GitHub App token
- `CLAUDE_CODE_OAUTH_TOKEN` not configured
- Incorrect permissions in `additional-permissions` input

**Solution:**
- Verify secrets are configured: `CLAUDE_APP_ID`, `CLAUDE_APP_PRIVATE_KEY`, `CLAUDE_CODE_OAUTH_TOKEN`
- Check that workflow uses `workflow_call` not direct event triggers

### Issue: Claude says it doesn't have tool access

**Symptoms:**
- Claude reports it cannot use `go`, `npm`, `git`, etc.
- Tool execution fails even though tools are installed

**Root cause:**
- `allowed-tools` input is empty or not passed through
- Usually caused by using direct event triggers instead of `workflow_call`
- Inputs are blank when workflow is triggered by events (not workflow_call)

**Solution:**
- Ensure workflows use `workflow_call` pattern:
  ```yaml
  # ✅ Correct - reusable workflow
  on:
    workflow_call:
      inputs:
        allowed-tools:
          type: string
          default: 'Bash(go:*),Bash(npm:*),Bash(git:*)'
  ```

  ```yaml
  # ❌ Wrong - direct event trigger
  on:
    issues:
      types: [opened]
    # This pattern leaves inputs blank!
  ```

- Use an internal wrapper workflow that listens for events and calls the reusable workflow:
  ```yaml
  # _claude.yml - internal wrapper
  on:
    issues:
      types: [opened]

  jobs:
    claude:
      uses: ./.github/workflows/claude.yml
      secrets: inherit
  ```

### Issue: Workflow doesn't trigger at all

**Symptoms:**
- @claude mention on issue/PR but no workflow runs

**Root cause:**
- Missing trigger conditions in internal wrapper workflow
- Workflow filter conditions too restrictive

**Solution:**
- Verify internal wrapper workflow (_claude.yml) has correct triggers
- Check that filter conditions properly detect @claude mentions
- Example:
  ```yaml
  jobs:
    claude:
      if: contains(github.event.issue.body, '@claude') || contains(github.event.comment.body, '@claude')
      uses: ./.github/workflows/claude.yml
      secrets: inherit
  ```

## Best Practices

### 1. Always use workflow_call for reusable workflows

Reusable workflows should ONLY use `workflow_call` trigger:

```yaml
# claude.yml
on:
  workflow_call:
    inputs:
      allowed-tools:
        type: string
        default: 'Bash(go:*),Bash(npm:*),...'
```

### 2. Create internal wrapper workflows for events

Internal workflows listen for events and call reusable workflows:

```yaml
# _claude.yml
on:
  issues:
    types: [opened]
  issue_comment:
    types: [created]

jobs:
  claude:
    uses: ./.github/workflows/claude.yml
    secrets: inherit
```

### 3. Always inherit secrets

```yaml
jobs:
  claude:
    uses: ./.github/workflows/claude.yml
    secrets: inherit  # Critical!
```

### 4. Test workflow changes

- Run `_test-claude-workflows.yml` after modifying Claude workflows
- Create test issues/PRs to verify end-to-end functionality
- Check workflow logs for permission or configuration errors

## Integration with PR Checks

The `_test-claude-workflows.yml` test runs automatically as a required check on PRs that modify:
- `.github/workflows/claude.yml`
- `.github/workflows/claude-code-review.yml`
- `.github/workflows/_test-claude-workflows.yml`

This ensures that structural issues are caught before merging changes that could break Claude functionality.

## Limitations

The automated test validates **workflow structure and configuration**, but cannot test:

- Actual Claude Code execution (requires CLAUDE_CODE_OAUTH_TOKEN)
- Comment update permissions (requires live GitHub events)
- Tool execution at runtime (requires actual Claude invocation)

For full validation, manual testing with real issues/PRs is necessary after structural validation passes.

## Contributing

When modifying Claude workflows:

1. Run the test workflow locally or via workflow_dispatch
2. Verify all checks pass
3. Test manually with a test issue or PR
4. Document any new configuration options
5. Update this guide if adding new validation checks
