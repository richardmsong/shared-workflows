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

## Important Notes

- When updating workflows, ensure changes are tested via the `_test-ci.yml` workflow
- Follow GitHub Actions best practices for reusable workflows
- Consider security implications when modifying workflow permissions or secrets handling
