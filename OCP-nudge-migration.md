# NudgeConfig Dry-Run Migration Report — `tekton-ecosystem-tenant`

**Cluster:** `kflux-prd-rh02`
**Namespace:** `tekton-ecosystem-tenant`
**Script:** upstream `nudge-migrate.sh` (MAX_NUDGES=360)
**Date:** 2026-09-15
**Mode:** DRY RUN — read-only, zero writes

---

## Summary

| Metric | Count |
|--------|-------|
| Raw nudge relationships found | 334 |
| Dangling references filtered | 24 |
| Valid relationships → NudgeConfig | **310** |
| NudgeConfig action | **CREATE** (none exists yet) |
| Cycles detected | None |
| Within 360 limit? | **Yes** (310/360) |

---

## Dangling References (24 — filtered, NOT migrated)

These point to components that no longer exist and will be skipped.

| Source Component | Application | Missing Target |
|---|---|---|
| `console-plugin-1-23-console-plugin-pf5` | `openshift-pipelines-core-1-23` | `operator-1-23-bundle-pf5` *(deleted)* |
| `multicluster-proxy-aae-0-1-multicluster-proxy-aae` | `multicluster-proxy-aae-0-1` | `operator-0-1-bundle` *(deleted)* |
| `operator-1-15-index-4-22` | `openshift-pipelines-index-4-22-1-15` | *(empty target)* |
| `operator-1-15-index-4-23` | `openshift-pipelines-index-4-23-1-15` | *(empty target)* |
| `operator-1-15-index-5-0` | `openshift-pipelines-index-5-0-1-15` | *(empty target)* |
| `operator-1-20-index-4-17` | `openshift-pipelines-index-4-17-1-20` | *(empty target)* |
| `operator-1-20-index-4-22` | `openshift-pipelines-index-4-22-1-20` | *(empty target)* |
| `operator-1-20-index-4-23` | `openshift-pipelines-index-4-23-1-20` | *(empty target)* |
| `operator-1-20-index-5-0` | `openshift-pipelines-index-5-0-1-20` | *(empty target)* |
| `operator-1-21-index-4-17` | `openshift-pipelines-index-4-17-1-21` | *(empty target)* |
| `operator-1-22-index-4-17` | `openshift-pipelines-index-4-17-1-22` | *(empty target)* |
| `operator-1-24-index-4-14` | `openshift-pipelines-index-4-14-1-24` | *(empty target)* |
| `operator-1-24-index-4-16` | `openshift-pipelines-index-4-16-1-24` | *(empty target)* |
| `operator-1-24-index-4-18` | `openshift-pipelines-index-4-18-1-24` | *(empty target)* |
| `operator-1-24-index-4-19` | `openshift-pipelines-index-4-19-1-24` | *(empty target)* |
| `operator-1-24-index-4-20` | `openshift-pipelines-index-4-20-1-24` | *(empty target)* |
| `operator-1-24-index-4-21` | `openshift-pipelines-index-4-21-1-24` | *(empty target)* |
| `operator-1-24-index-4-22` | `openshift-pipelines-index-4-22-1-24` | *(empty target)* |
| `operator-1-24-index-4-23` | `openshift-pipelines-index-4-23-1-24` | *(empty target)* |
| `operator-1-24-index-5-0` | `openshift-pipelines-index-5-0-1-24` | *(empty target)* |
| `pipelines-multikueue-plugin-0-1-controller` | `multicluster-pipelines-core-0-1` | `operator-0-1-bundle` *(deleted)* |
| `syncer-service-0-1-syncer-service` | `syncer-service-0-1` | `operator-0-1-bundle` *(deleted)* |
| `tekton-caches-0-3-cache` | `tekton-caches-0-3` | `operator-0-3-bundle` *(deleted)* |
| `tekton-caches-0-4-cache` | `openshift-pipelines-caches-0-4` | `operator-0-4-bundle` *(deleted)* |

---

## Valid Nudge Relationships (310 — will be migrated)

All sub-components per release branch converge onto their `operator-*-bundle`.
Each bundle then nudges its single index image (`operator-*-index-4-20`).
Releases covered: `1-15`, `1-20`, `1-21`, `1-22`, `1-23`, `1-24`, `next`, `nightly`.

---

### `openshift-pipelines-core-1-15` — 34 relationships

| Source Component | → Target |
|---|---|
| `console-plugin-1-15-console-plugin` | `operator-1-15-bundle` |
| `console-plugin-pf5-1-15-console-plugin` | `operator-1-15-bundle` |
| `manual-approval-gate-1-15-controller` | `operator-1-15-bundle` |
| `manual-approval-gate-1-15-webhook` | `operator-1-15-bundle` |
| `opc-1-15-opc` | `operator-1-15-bundle` |
| `operator-1-15-operator` | `operator-1-15-bundle` |
| `operator-1-15-proxy` | `operator-1-15-bundle` |
| `operator-1-15-webhook` | `operator-1-15-bundle` |
| `pipelines-as-code-1-15-cli` | `operator-1-15-bundle` |
| `pipelines-as-code-1-15-controller` | `operator-1-15-bundle` |
| `pipelines-as-code-1-15-watcher` | `operator-1-15-bundle` |
| `pipelines-as-code-1-15-webhook` | `operator-1-15-bundle` |
| `serve-tkn-cli-1-15-serve-tkn-cli` | `operator-1-15-bundle` |
| `tektoncd-chains-1-15-controller` | `operator-1-15-bundle` |
| `tektoncd-cli-1-15-tkn` | `operator-1-15-bundle` |
| `tektoncd-git-clone-1-15-git-init` | `operator-1-15-bundle` |
| `tektoncd-hub-1-15-api` | `operator-1-15-bundle` |
| `tektoncd-hub-1-15-db-migration` | `operator-1-15-bundle` |
| `tektoncd-hub-1-15-ui` | `operator-1-15-bundle` |
| `tektoncd-pipeline-1-15-controller` | `operator-1-15-bundle` |
| `tektoncd-pipeline-1-15-entrypoint` | `operator-1-15-bundle` |
| `tektoncd-pipeline-1-15-events` | `operator-1-15-bundle` |
| `tektoncd-pipeline-1-15-nop` | `operator-1-15-bundle` |
| `tektoncd-pipeline-1-15-resolvers` | `operator-1-15-bundle` |
| `tektoncd-pipeline-1-15-sidecarlogresults` | `operator-1-15-bundle` |
| `tektoncd-pipeline-1-15-webhook` | `operator-1-15-bundle` |
| `tektoncd-pipeline-1-15-workingdirinit` | `operator-1-15-bundle` |
| `tektoncd-results-1-15-api` | `operator-1-15-bundle` |
| `tektoncd-results-1-15-retention-policy-agent` | `operator-1-15-bundle` |
| `tektoncd-results-1-15-watcher` | `operator-1-15-bundle` |
| `tektoncd-triggers-1-15-controller` | `operator-1-15-bundle` |
| `tektoncd-triggers-1-15-core-interceptors` | `operator-1-15-bundle` |
| `tektoncd-triggers-1-15-eventlistenersink` | `operator-1-15-bundle` |
| `tektoncd-triggers-1-15-webhook` | `operator-1-15-bundle` |

**`openshift-pipelines-bundle-1-15`:**

| Source Component | → Target |
|---|---|
| `operator-1-15-bundle` | `operator-1-15-index-4-20` |

---

### `openshift-pipelines-core-1-20` — 37 relationships

| Source Component | → Target |
|---|---|
| `console-plugin-1-20-console-plugin` | `operator-1-20-bundle` |
| `console-plugin-pf5-1-20-console-plugin` | `operator-1-20-bundle` |
| `manual-approval-gate-1-20-controller` | `operator-1-20-bundle` |
| `manual-approval-gate-1-20-webhook` | `operator-1-20-bundle` |
| `opc-1-20-opc` | `operator-1-20-bundle` |
| `operator-1-20-operator` | `operator-1-20-bundle` |
| `operator-1-20-proxy` | `operator-1-20-bundle` |
| `operator-1-20-webhook` | `operator-1-20-bundle` |
| `pipelines-as-code-1-20-cli` | `operator-1-20-bundle` |
| `pipelines-as-code-1-20-controller` | `operator-1-20-bundle` |
| `pipelines-as-code-1-20-watcher` | `operator-1-20-bundle` |
| `pipelines-as-code-1-20-webhook` | `operator-1-20-bundle` |
| `serve-tkn-cli-1-20-serve-tkn-cli` | `operator-1-20-bundle` |
| `tekton-caches-1-20-cache` | `operator-1-20-bundle` |
| `tektoncd-chains-1-20-controller` | `operator-1-20-bundle` |
| `tektoncd-cli-1-20-tkn` | `operator-1-20-bundle` |
| `tektoncd-git-clone-1-20-git-init` | `operator-1-20-bundle` |
| `tektoncd-hub-1-20-api` | `operator-1-20-bundle` |
| `tektoncd-hub-1-20-db-migration` | `operator-1-20-bundle` |
| `tektoncd-hub-1-20-ui` | `operator-1-20-bundle` |
| `tektoncd-pipeline-1-20-controller` | `operator-1-20-bundle` |
| `tektoncd-pipeline-1-20-entrypoint` | `operator-1-20-bundle` |
| `tektoncd-pipeline-1-20-events` | `operator-1-20-bundle` |
| `tektoncd-pipeline-1-20-nop` | `operator-1-20-bundle` |
| `tektoncd-pipeline-1-20-resolvers` | `operator-1-20-bundle` |
| `tektoncd-pipeline-1-20-sidecarlogresults` | `operator-1-20-bundle` |
| `tektoncd-pipeline-1-20-webhook` | `operator-1-20-bundle` |
| `tektoncd-pipeline-1-20-workingdirinit` | `operator-1-20-bundle` |
| `tektoncd-pruner-1-20-controller` | `operator-1-20-bundle` |
| `tektoncd-pruner-1-20-webhook` | `operator-1-20-bundle` |
| `tektoncd-results-1-20-api` | `operator-1-20-bundle` |
| `tektoncd-results-1-20-retention-policy-agent` | `operator-1-20-bundle` |
| `tektoncd-results-1-20-watcher` | `operator-1-20-bundle` |
| `tektoncd-triggers-1-20-controller` | `operator-1-20-bundle` |
| `tektoncd-triggers-1-20-core-interceptors` | `operator-1-20-bundle` |
| `tektoncd-triggers-1-20-eventlistenersink` | `operator-1-20-bundle` |
| `tektoncd-triggers-1-20-webhook` | `operator-1-20-bundle` |

**`openshift-pipelines-bundle-1-20`:**

| Source Component | → Target |
|---|---|
| `operator-1-20-bundle` | `operator-1-20-index-4-20` |

---

### `openshift-pipelines-core-1-21` — 37 relationships

| Source Component | → Target |
|---|---|
| `console-plugin-1-21-console-plugin` | `operator-1-21-bundle` |
| `console-plugin-pf5-1-21-console-plugin` | `operator-1-21-bundle` |
| `manual-approval-gate-1-21-controller` | `operator-1-21-bundle` |
| `manual-approval-gate-1-21-webhook` | `operator-1-21-bundle` |
| `opc-1-21-opc` | `operator-1-21-bundle` |
| `operator-1-21-operator` | `operator-1-21-bundle` |
| `operator-1-21-proxy` | `operator-1-21-bundle` |
| `operator-1-21-webhook` | `operator-1-21-bundle` |
| `pipelines-as-code-1-21-cli` | `operator-1-21-bundle` |
| `pipelines-as-code-1-21-controller` | `operator-1-21-bundle` |
| `pipelines-as-code-1-21-watcher` | `operator-1-21-bundle` |
| `pipelines-as-code-1-21-webhook` | `operator-1-21-bundle` |
| `serve-tkn-cli-1-21-serve-tkn-cli` | `operator-1-21-bundle` |
| `tekton-caches-1-21-cache` | `operator-1-21-bundle` |
| `tektoncd-chains-1-21-controller` | `operator-1-21-bundle` |
| `tektoncd-cli-1-21-tkn` | `operator-1-21-bundle` |
| `tektoncd-git-clone-1-21-git-init` | `operator-1-21-bundle` |
| `tektoncd-hub-1-21-api` | `operator-1-21-bundle` |
| `tektoncd-hub-1-21-db-migration` | `operator-1-21-bundle` |
| `tektoncd-hub-1-21-ui` | `operator-1-21-bundle` |
| `tektoncd-pipeline-1-21-controller` | `operator-1-21-bundle` |
| `tektoncd-pipeline-1-21-entrypoint` | `operator-1-21-bundle` |
| `tektoncd-pipeline-1-21-events` | `operator-1-21-bundle` |
| `tektoncd-pipeline-1-21-nop` | `operator-1-21-bundle` |
| `tektoncd-pipeline-1-21-resolvers` | `operator-1-21-bundle` |
| `tektoncd-pipeline-1-21-sidecarlogresults` | `operator-1-21-bundle` |
| `tektoncd-pipeline-1-21-webhook` | `operator-1-21-bundle` |
| `tektoncd-pipeline-1-21-workingdirinit` | `operator-1-21-bundle` |
| `tektoncd-pruner-1-21-controller` | `operator-1-21-bundle` |
| `tektoncd-pruner-1-21-webhook` | `operator-1-21-bundle` |
| `tektoncd-results-1-21-api` | `operator-1-21-bundle` |
| `tektoncd-results-1-21-retention-policy-agent` | `operator-1-21-bundle` |
| `tektoncd-results-1-21-watcher` | `operator-1-21-bundle` |
| `tektoncd-triggers-1-21-controller` | `operator-1-21-bundle` |
| `tektoncd-triggers-1-21-core-interceptors` | `operator-1-21-bundle` |
| `tektoncd-triggers-1-21-eventlistenersink` | `operator-1-21-bundle` |
| `tektoncd-triggers-1-21-webhook` | `operator-1-21-bundle` |

**`openshift-pipelines-bundle-1-21`:**

| Source Component | → Target |
|---|---|
| `operator-1-21-bundle` | `operator-1-21-index-4-20` |

---

### `openshift-pipelines-core-1-22` — 40 relationships

| Source Component | → Target |
|---|---|
| `console-plugin-1-22-console-plugin` | `operator-1-22-bundle` |
| `console-plugin-pf5-1-22-console-plugin` | `operator-1-22-bundle` |
| `manual-approval-gate-1-22-controller` | `operator-1-22-bundle` |
| `manual-approval-gate-1-22-webhook` | `operator-1-22-bundle` |
| `multicluster-proxy-aae-1-22-multicluster-proxy-aae` | `operator-1-22-bundle` |
| `opc-1-22-opc` | `operator-1-22-bundle` |
| `operator-1-22-operator` | `operator-1-22-bundle` |
| `operator-1-22-proxy` | `operator-1-22-bundle` |
| `operator-1-22-webhook` | `operator-1-22-bundle` |
| `pipelines-as-code-1-22-cli` | `operator-1-22-bundle` |
| `pipelines-as-code-1-22-controller` | `operator-1-22-bundle` |
| `pipelines-as-code-1-22-watcher` | `operator-1-22-bundle` |
| `pipelines-as-code-1-22-webhook` | `operator-1-22-bundle` |
| `serve-tkn-cli-1-22-serve-tkn-cli` | `operator-1-22-bundle` |
| `syncer-service-1-22-syncer-service` | `operator-1-22-bundle` |
| `tekton-caches-1-22-cache` | `operator-1-22-bundle` |
| `tekton-kueue-1-22-scheduler` | `operator-1-22-bundle` |
| `tektoncd-chains-1-22-controller` | `operator-1-22-bundle` |
| `tektoncd-cli-1-22-tkn` | `operator-1-22-bundle` |
| `tektoncd-git-clone-1-22-git-init` | `operator-1-22-bundle` |
| `tektoncd-hub-1-22-api` | `operator-1-22-bundle` |
| `tektoncd-hub-1-22-db-migration` | `operator-1-22-bundle` |
| `tektoncd-hub-1-22-ui` | `operator-1-22-bundle` |
| `tektoncd-pipeline-1-22-controller` | `operator-1-22-bundle` |
| `tektoncd-pipeline-1-22-entrypoint` | `operator-1-22-bundle` |
| `tektoncd-pipeline-1-22-events` | `operator-1-22-bundle` |
| `tektoncd-pipeline-1-22-nop` | `operator-1-22-bundle` |
| `tektoncd-pipeline-1-22-resolvers` | `operator-1-22-bundle` |
| `tektoncd-pipeline-1-22-sidecarlogresults` | `operator-1-22-bundle` |
| `tektoncd-pipeline-1-22-webhook` | `operator-1-22-bundle` |
| `tektoncd-pipeline-1-22-workingdirinit` | `operator-1-22-bundle` |
| `tektoncd-pruner-1-22-controller` | `operator-1-22-bundle` |
| `tektoncd-pruner-1-22-webhook` | `operator-1-22-bundle` |
| `tektoncd-results-1-22-api` | `operator-1-22-bundle` |
| `tektoncd-results-1-22-retention-policy-agent` | `operator-1-22-bundle` |
| `tektoncd-results-1-22-watcher` | `operator-1-22-bundle` |
| `tektoncd-triggers-1-22-controller` | `operator-1-22-bundle` |
| `tektoncd-triggers-1-22-core-interceptors` | `operator-1-22-bundle` |
| `tektoncd-triggers-1-22-eventlistenersink` | `operator-1-22-bundle` |
| `tektoncd-triggers-1-22-webhook` | `operator-1-22-bundle` |

**`openshift-pipelines-bundle-1-22`:**

| Source Component | → Target |
|---|---|
| `operator-1-22-bundle` | `operator-1-22-index-4-20` |

---

### `openshift-pipelines-core-1-23` — 41 relationships

> Note: `console-plugin-1-23-console-plugin-pf5` had 2 nudges; the target `operator-1-23-bundle-pf5` is dangling and was filtered — only its valid nudge to `operator-1-23-bundle` is migrated.

| Source Component | → Target |
|---|---|
| `console-plugin-1-23-console-plugin` | `operator-1-23-bundle` |
| `console-plugin-1-23-console-plugin-pf5` | `operator-1-23-bundle` |
| `console-plugin-pf5-1-23-console-plugin` | `operator-1-23-bundle` |
| `manual-approval-gate-1-23-controller` | `operator-1-23-bundle` |
| `manual-approval-gate-1-23-webhook` | `operator-1-23-bundle` |
| `multicluster-proxy-aae-1-23-multicluster-proxy-aae` | `operator-1-23-bundle` |
| `opc-1-23-opc` | `operator-1-23-bundle` |
| `operator-1-23-operator` | `operator-1-23-bundle` |
| `operator-1-23-proxy` | `operator-1-23-bundle` |
| `operator-1-23-webhook` | `operator-1-23-bundle` |
| `pipelines-as-code-1-23-cli` | `operator-1-23-bundle` |
| `pipelines-as-code-1-23-controller` | `operator-1-23-bundle` |
| `pipelines-as-code-1-23-watcher` | `operator-1-23-bundle` |
| `pipelines-as-code-1-23-webhook` | `operator-1-23-bundle` |
| `serve-tkn-cli-1-23-serve-tkn-cli` | `operator-1-23-bundle` |
| `syncer-service-1-23-syncer-service` | `operator-1-23-bundle` |
| `tekton-caches-1-23-cache` | `operator-1-23-bundle` |
| `tekton-kueue-1-23-scheduler` | `operator-1-23-bundle` |
| `tektoncd-chains-1-23-controller` | `operator-1-23-bundle` |
| `tektoncd-cli-1-23-tkn` | `operator-1-23-bundle` |
| `tektoncd-git-clone-1-23-git-init` | `operator-1-23-bundle` |
| `tektoncd-hub-1-23-api` | `operator-1-23-bundle` |
| `tektoncd-hub-1-23-db-migration` | `operator-1-23-bundle` |
| `tektoncd-hub-1-23-ui` | `operator-1-23-bundle` |
| `tektoncd-pipeline-1-23-controller` | `operator-1-23-bundle` |
| `tektoncd-pipeline-1-23-entrypoint` | `operator-1-23-bundle` |
| `tektoncd-pipeline-1-23-events` | `operator-1-23-bundle` |
| `tektoncd-pipeline-1-23-nop` | `operator-1-23-bundle` |
| `tektoncd-pipeline-1-23-resolvers` | `operator-1-23-bundle` |
| `tektoncd-pipeline-1-23-sidecarlogresults` | `operator-1-23-bundle` |
| `tektoncd-pipeline-1-23-webhook` | `operator-1-23-bundle` |
| `tektoncd-pipeline-1-23-workingdirinit` | `operator-1-23-bundle` |
| `tektoncd-pruner-1-23-controller` | `operator-1-23-bundle` |
| `tektoncd-pruner-1-23-webhook` | `operator-1-23-bundle` |
| `tektoncd-results-1-23-api` | `operator-1-23-bundle` |
| `tektoncd-results-1-23-retention-policy-agent` | `operator-1-23-bundle` |
| `tektoncd-results-1-23-watcher` | `operator-1-23-bundle` |
| `tektoncd-triggers-1-23-controller` | `operator-1-23-bundle` |
| `tektoncd-triggers-1-23-core-interceptors` | `operator-1-23-bundle` |
| `tektoncd-triggers-1-23-eventlistenersink` | `operator-1-23-bundle` |
| `tektoncd-triggers-1-23-webhook` | `operator-1-23-bundle` |

**`openshift-pipelines-bundle-1-23`:**

| Source Component | → Target |
|---|---|
| `operator-1-23-bundle` | `operator-1-23-index-4-20` |

---

### `openshift-pipelines-core-1-24` — 37 relationships

| Source Component | → Target |
|---|---|
| `console-plugin-1-24-console-plugin` | `operator-1-24-bundle` |
| `console-plugin-pf5-1-24-console-plugin` | `operator-1-24-bundle` |
| `manual-approval-gate-1-24-controller` | `operator-1-24-bundle` |
| `manual-approval-gate-1-24-webhook` | `operator-1-24-bundle` |
| `multicluster-proxy-aae-1-24-multicluster-proxy-aae` | `operator-1-24-bundle` |
| `opc-1-24-opc` | `operator-1-24-bundle` |
| `operator-1-24-operator` | `operator-1-24-bundle` |
| `operator-1-24-proxy` | `operator-1-24-bundle` |
| `operator-1-24-webhook` | `operator-1-24-bundle` |
| `pipelines-as-code-1-24-cli` | `operator-1-24-bundle` |
| `pipelines-as-code-1-24-controller` | `operator-1-24-bundle` |
| `pipelines-as-code-1-24-watcher` | `operator-1-24-bundle` |
| `pipelines-as-code-1-24-webhook` | `operator-1-24-bundle` |
| `serve-tkn-cli-1-24-serve-tkn-cli` | `operator-1-24-bundle` |
| `syncer-service-1-24-syncer-service` | `operator-1-24-bundle` |
| `tekton-caches-1-24-cache` | `operator-1-24-bundle` |
| `tekton-kueue-1-24-scheduler` | `operator-1-24-bundle` |
| `tektoncd-chains-1-24-controller` | `operator-1-24-bundle` |
| `tektoncd-cli-1-24-tkn` | `operator-1-24-bundle` |
| `tektoncd-git-clone-1-24-git-init` | `operator-1-24-bundle` |
| `tektoncd-pipeline-1-24-controller` | `operator-1-24-bundle` |
| `tektoncd-pipeline-1-24-entrypoint` | `operator-1-24-bundle` |
| `tektoncd-pipeline-1-24-events` | `operator-1-24-bundle` |
| `tektoncd-pipeline-1-24-nop` | `operator-1-24-bundle` |
| `tektoncd-pipeline-1-24-resolvers` | `operator-1-24-bundle` |
| `tektoncd-pipeline-1-24-sidecarlogresults` | `operator-1-24-bundle` |
| `tektoncd-pipeline-1-24-webhook` | `operator-1-24-bundle` |
| `tektoncd-pipeline-1-24-workingdirinit` | `operator-1-24-bundle` |
| `tektoncd-pruner-1-24-controller` | `operator-1-24-bundle` |
| `tektoncd-pruner-1-24-webhook` | `operator-1-24-bundle` |
| `tektoncd-results-1-24-api` | `operator-1-24-bundle` |
| `tektoncd-results-1-24-retention-policy-agent` | `operator-1-24-bundle` |
| `tektoncd-results-1-24-watcher` | `operator-1-24-bundle` |
| `tektoncd-triggers-1-24-controller` | `operator-1-24-bundle` |
| `tektoncd-triggers-1-24-core-interceptors` | `operator-1-24-bundle` |
| `tektoncd-triggers-1-24-eventlistenersink` | `operator-1-24-bundle` |
| `tektoncd-triggers-1-24-webhook` | `operator-1-24-bundle` |

**`openshift-pipelines-bundle-1-24`:**

| Source Component | → Target |
|---|---|
| `operator-1-24-bundle` | `operator-1-24-index-4-20` |

---

### `openshift-pipelines-core-next` — 38 relationships

| Source Component | → Target |
|---|---|
| `console-plugin-next-console-plugin` | `operator-next-bundle` |
| `console-plugin-pf5-next-console-plugin` | `operator-next-bundle` |
| `manual-approval-gate-next-controller` | `operator-next-bundle` |
| `manual-approval-gate-next-webhook` | `operator-next-bundle` |
| `multicluster-proxy-aae-next-multicluster-proxy-aae` | `operator-next-bundle` |
| `opc-next-opc` | `operator-next-bundle` |
| `operator-next-operator` | `operator-next-bundle` |
| `operator-next-proxy` | `operator-next-bundle` |
| `operator-next-webhook` | `operator-next-bundle` |
| `pipelines-as-code-next-cli` | `operator-next-bundle` |
| `pipelines-as-code-next-controller` | `operator-next-bundle` |
| `pipelines-as-code-next-watcher` | `operator-next-bundle` |
| `pipelines-as-code-next-webhook` | `operator-next-bundle` |
| `pipelines-multikueue-plugin-next-multikueue-plugin` | `operator-next-bundle` |
| `serve-tkn-cli-next-serve-tkn-cli` | `operator-next-bundle` |
| `syncer-service-next-syncer-service` | `operator-next-bundle` |
| `tekton-caches-next-cache` | `operator-next-bundle` |
| `tekton-kueue-next-scheduler` | `operator-next-bundle` |
| `tektoncd-chains-next-controller` | `operator-next-bundle` |
| `tektoncd-cli-next-tkn` | `operator-next-bundle` |
| `tektoncd-git-clone-next-git-init` | `operator-next-bundle` |
| `tektoncd-pipeline-next-controller` | `operator-next-bundle` |
| `tektoncd-pipeline-next-entrypoint` | `operator-next-bundle` |
| `tektoncd-pipeline-next-events` | `operator-next-bundle` |
| `tektoncd-pipeline-next-nop` | `operator-next-bundle` |
| `tektoncd-pipeline-next-resolvers` | `operator-next-bundle` |
| `tektoncd-pipeline-next-sidecarlogresults` | `operator-next-bundle` |
| `tektoncd-pipeline-next-webhook` | `operator-next-bundle` |
| `tektoncd-pipeline-next-workingdirinit` | `operator-next-bundle` |
| `tektoncd-pruner-next-controller` | `operator-next-bundle` |
| `tektoncd-pruner-next-webhook` | `operator-next-bundle` |
| `tektoncd-results-next-api` | `operator-next-bundle` |
| `tektoncd-results-next-retention-policy-agent` | `operator-next-bundle` |
| `tektoncd-results-next-watcher` | `operator-next-bundle` |
| `tektoncd-triggers-next-controller` | `operator-next-bundle` |
| `tektoncd-triggers-next-core-interceptors` | `operator-next-bundle` |
| `tektoncd-triggers-next-eventlistenersink` | `operator-next-bundle` |
| `tektoncd-triggers-next-webhook` | `operator-next-bundle` |

**`openshift-pipelines-bundle-next`:**

| Source Component | → Target |
|---|---|
| `operator-next-bundle` | `operator-next-index-4-20` |

---

### `openshift-pipelines-core-nightly` — 38 relationships

| Source Component | → Target |
|---|---|
| `console-plugin-nightly-console-plugin` | `operator-nightly-bundle` |
| `console-plugin-pf5-nightly-console-plugin` | `operator-nightly-bundle` |
| `manual-approval-gate-nightly-controller` | `operator-nightly-bundle` |
| `manual-approval-gate-nightly-webhook` | `operator-nightly-bundle` |
| `multicluster-proxy-aae-nightly-multicluster-proxy-aae` | `operator-nightly-bundle` |
| `opc-nightly-opc` | `operator-nightly-bundle` |
| `operator-nightly-operator` | `operator-nightly-bundle` |
| `operator-nightly-proxy` | `operator-nightly-bundle` |
| `operator-nightly-webhook` | `operator-nightly-bundle` |
| `pipelines-as-code-nightly-cli` | `operator-nightly-bundle` |
| `pipelines-as-code-nightly-controller` | `operator-nightly-bundle` |
| `pipelines-as-code-nightly-watcher` | `operator-nightly-bundle` |
| `pipelines-as-code-nightly-webhook` | `operator-nightly-bundle` |
| `pipelines-multikueue-plugin-nightly-multikueue-plugin` | `operator-nightly-bundle` |
| `serve-tkn-cli-nightly-serve-tkn-cli` | `operator-nightly-bundle` |
| `syncer-service-nightly-syncer-service` | `operator-nightly-bundle` |
| `tekton-caches-nightly-cache` | `operator-nightly-bundle` |
| `tekton-kueue-nightly-scheduler` | `operator-nightly-bundle` |
| `tektoncd-chains-nightly-controller` | `operator-nightly-bundle` |
| `tektoncd-cli-nightly-tkn` | `operator-nightly-bundle` |
| `tektoncd-git-clone-nightly-git-init` | `operator-nightly-bundle` |
| `tektoncd-pipeline-nightly-controller` | `operator-nightly-bundle` |
| `tektoncd-pipeline-nightly-entrypoint` | `operator-nightly-bundle` |
| `tektoncd-pipeline-nightly-events` | `operator-nightly-bundle` |
| `tektoncd-pipeline-nightly-nop` | `operator-nightly-bundle` |
| `tektoncd-pipeline-nightly-resolvers` | `operator-nightly-bundle` |
| `tektoncd-pipeline-nightly-sidecarlogresults` | `operator-nightly-bundle` |
| `tektoncd-pipeline-nightly-webhook` | `operator-nightly-bundle` |
| `tektoncd-pipeline-nightly-workingdirinit` | `operator-nightly-bundle` |
| `tektoncd-pruner-nightly-controller` | `operator-nightly-bundle` |
| `tektoncd-pruner-nightly-webhook` | `operator-nightly-bundle` |
| `tektoncd-results-nightly-api` | `operator-nightly-bundle` |
| `tektoncd-results-nightly-retention-policy-agent` | `operator-nightly-bundle` |
| `tektoncd-results-nightly-watcher` | `operator-nightly-bundle` |
| `tektoncd-triggers-nightly-controller` | `operator-nightly-bundle` |
| `tektoncd-triggers-nightly-core-interceptors` | `operator-nightly-bundle` |
| `tektoncd-triggers-nightly-eventlistenersink` | `operator-nightly-bundle` |
| `tektoncd-triggers-nightly-webhook` | `operator-nightly-bundle` |

**`openshift-pipelines-bundle-nightly`:**

| Source Component | → Target |
|---|---|
| `operator-nightly-bundle` | `operator-nightly-index-4-20` |

---

## Conclusion

The namespace is **ready to migrate**. The graph is a clean star+chain pattern per release branch:
all sub-components converge onto their `operator-*-bundle`, and each bundle nudges one index image.

To apply (when ready):

```bash
bash nudge-migrate.sh tekton-ecosystem-tenant
```

Rollback if needed:

```bash
kubectl delete nudgeconfig nudge-config -n tekton-ecosystem-tenant
```
