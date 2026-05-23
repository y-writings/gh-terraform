# Explicit Token Config Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace token enable booleans with explicit optional token configuration objects that include all 1Password and GitHub output names.

**Architecture:** The root module owns reusable token presets in `locals` and passes either a complete object or `null` to the `release_please` module. The module no longer has token defaults or fixed token metadata; it only consumes complete objects and uses object presence to decide whether to create optional resources.

**Tech Stack:** Terraform `>= 1.7.0`, GitHub provider, 1Password provider.

---

### Task 1: Root Token Presets And Module Inputs

**Files:**
- Modify: `terraform/work-repositories/locals.tf`
- Modify: `terraform/work-repositories/main.tf`

- [ ] **Step 1: Update repository locals to use complete token objects**

Replace boolean flags with `token_presets` and `merge(...)` calls:

```hcl
locals {
  token_presets = {
    metrics = {
      item_title  = "metrics-token"
      section     = "info"
      field       = "token"
      secret_name = "METRICS_TOKEN"
    }

    pull_request_creator = {
      item_title              = "pull-request-creator-bot"
      app_id_section          = "info"
      app_id_field            = "app_id"
      private_key_secret_name = "PULL_REQUEST_CREATOR_APP_PRIVATE_KEY"
      app_id_variable_name    = "PULL_REQUEST_CREATOR_APP_ID"
    }
  }

  repositories = {
    repo_fe83b6f2 = {
      name = "y-writings"
      metrics_token = merge(local.token_presets.metrics, {
        vault_name = "dev"
      })
    }
    repo_6e7bb53d = {
      name = "calver-beacon-action"
      pull_request_creator_token = merge(local.token_presets.pull_request_creator, {
        vault_name = "dev"
      })
    }
  }
}
```

- [ ] **Step 2: Pass optional objects to the module**

Change the `module "release_please"` inputs to:

```hcl
  metrics_token        = try(each.value.metrics_token, null)
  pull_request_creator_token = try(each.value.pull_request_creator_token, null)
```

### Task 2: Module Variables And Resource Logic

**Files:**
- Modify: `terraform/modules/release_please/variables.tf`
- Modify: `terraform/modules/release_please/main.tf`

- [ ] **Step 1: Replace boolean variables with nullable complete objects**

Define `metrics_token` and `pull_request_creator_token` as nullable object variables with no default.

- [ ] **Step 2: Remove module token metadata defaults**

Remove fixed `vault_name`, release-please item metadata, and metrics token metadata from `locals`.

- [ ] **Step 3: Gate optional data sources and resources by object presence**

Use `count = var.pull_request_creator_token != null ? 1 : 0` for pull request creator 1Password data and GitHub resources. Use `count = var.metrics_token != null ? 1 : 0` for metrics 1Password data and GitHub secret.

- [ ] **Step 4: Read values from the supplied objects**

Use `var.pull_request_creator_token.vault_name`, `var.pull_request_creator_token.item_title`, `var.pull_request_creator_token.app_id_section`, `var.pull_request_creator_token.app_id_field`, `var.pull_request_creator_token.private_key_secret_name`, `var.pull_request_creator_token.app_id_variable_name`, and the equivalent `metrics_token` fields.

### Task 3: Docs And Verification

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update README examples and behavior notes**

Replace `enable_metrics_token` and `enable_pull_request_creator_token` examples with `metrics_token = merge(...)` and `pull_request_creator_token = merge(...)` examples. Update behavior text to say optional token objects control creation.

- [ ] **Step 2: Format Terraform**

Run: `terraform fmt -recursive terraform`

Expected: command exits 0.

- [ ] **Step 3: Validate the work root module**

Run: `terraform -chdir=terraform/work-repositories validate`

Expected: `Success! The configuration is valid.`

- [ ] **Step 4: Inspect diff**

Run: `git diff -- terraform/work-repositories terraform/modules/release_please README.md`

Expected: booleans are removed, complete token objects are used, and no unrelated files are changed.
