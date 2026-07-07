# nudge-migrate

One-time migration script that reads `build-nudges-ref` data from Component CRs
and creates or merges NudgeConfig CRDs per namespace.

This is part of **ADR-0067 Phase 2** — moving nudge relationship storage from
`Component.spec.build-nudges-ref` (build-service) to the `NudgeConfig` singleton
CRD (integration-service).

## Prerequisites

- The **NudgeConfig CRD** must be deployed on the target cluster
  (STONEINTG-1659, STONEINTG-1660).
- The **build-service skip patch** (STONEINTG-1672) must be deployed before
  running a live (non-dry-run) migration. Otherwise both services fire nudges
  simultaneously.
- The runner (user or service account) needs RBAC access to:
  - `list` Namespaces cluster-wide
  - `list` Components in all tenant namespaces
  - `get`, `create`, `update` NudgeConfigs in all tenant namespaces

## Usage

```bash
# Build
go build -o bin/nudge-migrate ./cmd/nudge-migrate/

# Dry-run all tenant namespaces (recommended first step)
./bin/nudge-migrate --dry-run

# Dry-run specific namespaces
./bin/nudge-migrate --dry-run tenant-foo tenant-bar

# Live migration — all tenant namespaces
./bin/nudge-migrate

# Live migration — specific namespaces
./bin/nudge-migrate tenant-foo tenant-bar
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--dry-run` | `false` | Print what would change without writing to the cluster |

Positional arguments are namespace names. If none are provided, the script
auto-discovers all tenant namespaces using label selectors.

## How It Works

### Namespace Discovery

When no namespaces are specified, the script discovers tenant namespaces using
three label queries (matching the `snapshotgc` pattern):

| Label | Value |
|-------|-------|
| `toolchain.dev.openshift.com/type` | `tenant` |
| `konflux.ci/type` | `user` |
| `konflux-ci.dev/type` | `tenant` |

Results are deduplicated and sorted alphabetically.

### Per-Namespace Logic

For each namespace:

1. **List all Components** in the namespace.
2. **Collect `build-nudges-ref`** entries from each Component.
3. **Filter out** invalid entries:
   - Self-nudges (`from == to`) are silently skipped.
   - Dangling references (target Component does not exist) are skipped with a
     warning log.
   - Duplicate `(from, to)` pairs are deduplicated.
4. If no valid relationships remain, **skip** the namespace.
5. **Validate the DAG** — run cycle detection via `dag.ValidateNudgeGraph()`.
   If a cycle is detected, the namespace is skipped with an error.
6. **Check cardinality** — if relationships exceed 256 (NudgeConfig CEL max),
   the namespace is skipped with an error.
7. **Check for existing NudgeConfig** (`nudge-config`) in the namespace:
   - **Not found** → create a new NudgeConfig with all relationships.
   - **Found** → merge: add entries from `build-nudges-ref` that don't already
     exist, preserving user-added entries.
8. On merge, re-validate the merged graph for cycles and cardinality.

### Merge Behavior

When a NudgeConfig already exists, the script **merges, never replaces**:

| Entry state | Action |
|-------------|--------|
| In `build-nudges-ref` but not in NudgeConfig | **Added** |
| Already in NudgeConfig (same `from`/`to`) | **Unchanged** |
| In NudgeConfig but not in `build-nudges-ref` (user-added) | **Preserved** |

This is safe because users may have started editing NudgeConfig directly after
an earlier script run.

### Migration Metadata

Every NudgeConfig created or updated by the script gets:

| Metadata | Key | Value |
|----------|-----|-------|
| Label | `nudging.konflux-ci.dev/owner` | `build-service` |
| Annotation | `nudging.konflux-ci.dev/migrated-from` | `build-nudges-ref` |

### Idempotency

Running the script twice on the same namespace is safe — the second run detects
that all relationships are already present and reports "skipped" with no changes.

## Output

### Dry-Run

```
[DRY RUN] Migration summary
  namespacesProcessed: 453
  created: 5
  updated: 2
  skipped: 440
  errors: 6
  totalRelationshipsMigrated: 23
```

### Live Run

```
Migration summary
  namespacesProcessed: 453
  created: 5
  updated: 2
  skipped: 440
  errors: 6
  totalRelationshipsMigrated: 23
```

The script exits with code 1 if any namespace encountered an error.

Per-namespace actions are logged at info level. Increase verbosity with `-v=1`
to see individual entries logged during dry-run creates.

## RBAC Requirements

The script must be run with a service account that has the following permissions.
Running with a personal user kubeconfig will fail with "forbidden" errors on
namespaces where the user lacks NudgeConfig access.

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: nudge-migrate
rules:
  - apiGroups: [""]
    resources: ["namespaces"]
    verbs: ["list"]
  - apiGroups: ["appstudio.redhat.com"]
    resources: ["components"]
    verbs: ["list"]
  - apiGroups: ["appstudio.redhat.com"]
    resources: ["nudgeconfigs"]
    verbs: ["get", "create", "update"]
```

**Note:** The `appstudio-admin-user-actions` role in infra-deployments must also
be updated to include `nudgeconfigs` verbs for tenant users to view/edit their
NudgeConfig after migration.

## Rollback

Deleting a NudgeConfig immediately reverts the namespace to build-service
nudging (assuming STONEINTG-1672 is deployed):

```bash
kubectl delete nudgeconfig nudge-config -n <namespace>
```

Build-service resumes nudging via `build-nudges-ref`. Re-run the migration
script to restore when ready.

## Cluster Test Results

Dry-run results from 2026-07-06 using a personal kubeconfig (limited RBAC):

### Production (`stone-prd-rh01`)

| Metric | Count |
|--------|-------|
| Namespaces discovered | 453 |
| Skipped (no nudges) | 407 |
| Errors (RBAC forbidden) | 46 |
| Dangling references | Multiple (filtered correctly) |

### Stage (`stone-stg-rh01`)

| Metric | Count |
|--------|-------|
| Namespaces discovered | 3,368 |
| Skipped (no nudges) | 3,358 |
| Errors (RBAC forbidden) | 10 |

Stage has far fewer nudge relationships than production (10 vs 46). The bulk of
stage namespaces are `test-rhtap-*` perfscale test namespaces with no components.

All RBAC errors are expected — the script needs a service account with
cluster-wide NudgeConfig access, not a personal user kubeconfig.

## Related Jira Issues

| Issue | Description |
|-------|-------------|
| STONEINTG-1682 | This migration script |
| STONEINTG-1659 | NudgeConfig CRD schema + CEL |
| STONEINTG-1660 | NudgeConfig webhook (cycle detection) |
| STONEINTG-1671 | Integration-service nudge orchestration |
| STONEINTG-1672 | Build-service skip when NudgeConfig exists |
