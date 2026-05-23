# Remove Applied Moved Block Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Remove the applied Terraform `moved` block after the state update has completed.

**Architecture:** The `moved` block exists in a dedicated Terraform file under the `work-repositories` root module. Since the file contains only the applied migration metadata, deleting the file is the smallest correct change and leaves runtime resources unchanged.

**Tech Stack:** Terraform `>= 1.7.0`, Git, GitHub CLI.

---

### Task 1: Remove Applied Moved Block

**Files:**
- Delete: `terraform/work-repositories/moved.tf`

- [ ] **Step 1: Delete the dedicated moved block file**

Remove `terraform/work-repositories/moved.tf`, which contains only this applied state migration metadata:

```hcl
moved {
  from = module.release_please
  to   = module.github_actions_credentials
}
```

- [ ] **Step 2: Format Terraform**

Run: `terraform fmt -recursive terraform`

Expected: command exits 0.

- [ ] **Step 3: Validate the work root module**

Run: `terraform -chdir=terraform/work-repositories validate`

Expected: `Success! The configuration is valid.`

- [ ] **Step 4: Inspect diff**

Run: `git diff -- terraform/work-repositories/moved.tf`

Expected: the only Terraform change is deletion of the applied `moved` block file.

- [ ] **Step 5: Commit cleanup**

```bash
git add terraform/work-repositories/moved.tf docs/superpowers/plans/2026-05-23-remove-applied-moved-block.md
git commit -m "chore(terraform): remove applied moved block"
```

- [ ] **Step 6: Push and create PR**

```bash
git push -u origin <branch-name>
gh pr create --title "chore(terraform): remove applied moved block" --body "## Summary\n- remove the applied Terraform moved block for github_actions_credentials\n\n## Tests\n- terraform fmt -recursive terraform\n- terraform -chdir=terraform/work-repositories validate"
```
