/*
Copyright 2026 Red Hat Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/go-logr/logr"
	applicationapiv1alpha1 "github.com/konflux-ci/application-api/api/v1alpha1"
	integrationv1beta2 "github.com/konflux-ci/integration-service/api/v1beta2"
	"github.com/konflux-ci/integration-service/pkg/dag"
	zap2 "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	core "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	scheme = runtime.NewScheme()
)

const (
	MigrationLabel           = "nudging.konflux-ci.dev/owner"
	MigrationLabelValue      = "build-service"
	MigrationAnnotation      = "nudging.konflux-ci.dev/migrated-from"
	MigrationAnnotationValue = "build-nudges-ref"
	MaxNudgesPerConfig       = 256
)

func init() {
	utilruntime.Must(applicationapiv1alpha1.AddToScheme(scheme))
	utilruntime.Must(integrationv1beta2.AddToScheme(scheme))
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
}

type migrationSummary struct {
	Namespace                    string
	ComponentsRead               int
	RelationshipsFound           int
	RelationshipsMigrated        int
	RelationshipsSkippedDangling int
	RelationshipsAlreadyPresent  int
	Action                       string // "created", "updated", "skipped", "dry-run:create", "dry-run:update", "dry-run:skip", "error"
	Errors                       []string
}

func main() {
	var dryRun bool
	flag.BoolVar(&dryRun, "dry-run", false, "Print what would change without writing to the cluster")

	opts := zap.Options{
		Development: false,
		TimeEncoder: zapcore.RFC3339TimeEncoder,
		ZapOpts:     []zap2.Option{zap2.WithCaller(true)},
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	logger := zap.New(zap.UseFlagOptions(&opts))
	namespaces := flag.Args()

	cl, err := client.New(config.GetConfigOrDie(), client.Options{Scheme: scheme})
	if err != nil {
		logger.Error(err, "Failed to create Kubernetes client")
		os.Exit(1)
	}

	ctx := context.Background()

	if len(namespaces) == 0 {
		logger.Info("No namespaces specified, discovering tenant namespaces")
		namespaces, err = getTenantNamespaces(ctx, cl, logger)
		if err != nil {
			logger.Error(err, "Failed to discover tenant namespaces")
			os.Exit(1)
		}
		logger.Info("Discovered tenant namespaces", "count", len(namespaces))
	}

	summaries := migrateNamespaces(ctx, cl, logger, namespaces, dryRun)
	printSummary(summaries, logger, dryRun)

	for _, s := range summaries {
		if len(s.Errors) > 0 {
			os.Exit(1)
		}
	}
}

func getTenantNamespaces(ctx context.Context, cl client.Client, logger logr.Logger) ([]string, error) {
	labelQueries := []struct {
		key   string
		value string
	}{
		{"toolchain.dev.openshift.com/type", "tenant"},
		{"konflux.ci/type", "user"},
		{"konflux-ci.dev/type", "tenant"},
	}

	seen := make(map[string]bool)
	var result []string

	for _, lq := range labelQueries {
		req, err := labels.NewRequirement(lq.key, selection.In, []string{lq.value})
		if err != nil {
			return nil, fmt.Errorf("building label requirement %s=%s: %w", lq.key, lq.value, err)
		}
		selector := labels.NewSelector().Add(*req)
		nsList := &core.NamespaceList{}
		if err := cl.List(ctx, nsList, &client.ListOptions{LabelSelector: selector}); err != nil {
			logger.Error(err, "Failed listing namespaces", "label", lq.key+"="+lq.value)
			return nil, err
		}
		for i := range nsList.Items {
			name := nsList.Items[i].Name
			if !seen[name] {
				seen[name] = true
				result = append(result, name)
			}
		}
	}

	sort.Strings(result)
	return result, nil
}

func migrateNamespaces(ctx context.Context, cl client.Client, logger logr.Logger, namespaces []string, dryRun bool) []migrationSummary {
	summaries := make([]migrationSummary, 0, len(namespaces))
	for _, ns := range namespaces {
		s := migrateNamespace(ctx, cl, logger.WithValues("namespace", ns), ns, dryRun)
		summaries = append(summaries, s)
	}
	return summaries
}

func migrateNamespace(ctx context.Context, cl client.Client, logger logr.Logger, namespace string, dryRun bool) migrationSummary {
	summary := migrationSummary{Namespace: namespace}

	componentList := &applicationapiv1alpha1.ComponentList{}
	if err := cl.List(ctx, componentList, client.InNamespace(namespace)); err != nil {
		logger.Error(err, "Failed to list Components")
		summary.Action = "error"
		summary.Errors = append(summary.Errors, fmt.Sprintf("list components: %v", err))
		return summary
	}
	summary.ComponentsRead = len(componentList.Items)

	existingComponents := make(map[string]bool, len(componentList.Items))
	for i := range componentList.Items {
		existingComponents[componentList.Items[i].Name] = true
	}

	var relationships []integrationv1beta2.NudgeRelationship
	seen := make(map[string]bool)

	for i := range componentList.Items {
		comp := &componentList.Items[i]
		for _, target := range comp.Spec.BuildNudgesRef {
			if comp.Name == target {
				logger.Info("Skipping self-nudge", "component", comp.Name)
				continue
			}
			if !existingComponents[target] {
				logger.Info("Skipping dangling reference (target component does not exist)",
					"from", comp.Name, "to", target)
				summary.RelationshipsSkippedDangling++
				continue
			}
			key := comp.Name + "->" + target
			if seen[key] {
				continue
			}
			seen[key] = true
			relationships = append(relationships, integrationv1beta2.NudgeRelationship{
				From: comp.Name,
				To:   target,
				Mode: integrationv1beta2.NudgeModeImmediate,
			})
		}
	}
	summary.RelationshipsFound = len(relationships)

	if len(relationships) == 0 {
		if dryRun {
			summary.Action = "dry-run:skip"
		} else {
			summary.Action = "skipped"
		}
		logger.Info("No nudge relationships found, skipping namespace")
		return summary
	}

	if err := dag.ValidateNudgeGraph(relationships); err != nil {
		logger.Error(err, "Cycle detected in nudge relationships, skipping namespace")
		summary.Action = "error"
		summary.Errors = append(summary.Errors, fmt.Sprintf("cycle detected: %v", err))
		return summary
	}

	if len(relationships) > MaxNudgesPerConfig {
		err := fmt.Errorf("found %d relationships, exceeds maximum %d", len(relationships), MaxNudgesPerConfig)
		logger.Error(err, "Too many nudge relationships, skipping namespace")
		summary.Action = "error"
		summary.Errors = append(summary.Errors, err.Error())
		return summary
	}

	existing := &integrationv1beta2.NudgeConfig{}
	err := cl.Get(ctx, client.ObjectKey{Namespace: namespace, Name: integrationv1beta2.NudgeConfigSingletonName}, existing)

	if apierrors.IsNotFound(err) {
		return handleCreate(ctx, cl, logger, namespace, relationships, dryRun, summary)
	}
	if err != nil {
		logger.Error(err, "Failed to get existing NudgeConfig")
		summary.Action = "error"
		summary.Errors = append(summary.Errors, fmt.Sprintf("get nudgeconfig: %v", err))
		return summary
	}

	return handleMerge(ctx, cl, logger, existing, relationships, dryRun, summary)
}

func handleCreate(ctx context.Context, cl client.Client, logger logr.Logger, namespace string, relationships []integrationv1beta2.NudgeRelationship, dryRun bool, summary migrationSummary) migrationSummary {
	nc := buildNudgeConfig(namespace, relationships)
	summary.RelationshipsMigrated = len(relationships)

	if dryRun {
		summary.Action = "dry-run:create"
		logger.Info("Would CREATE NudgeConfig",
			"entries", len(relationships))
		for _, r := range relationships {
			logger.V(1).Info("  entry", "from", r.From, "to", r.To, "mode", r.Mode)
		}
		return summary
	}

	if err := cl.Create(ctx, nc); err != nil {
		logger.Error(err, "Failed to create NudgeConfig")
		summary.Action = "error"
		summary.Errors = append(summary.Errors, fmt.Sprintf("create nudgeconfig: %v", err))
		return summary
	}

	summary.Action = "created"
	logger.Info("Created NudgeConfig", "entries", len(relationships))
	return summary
}

func handleMerge(ctx context.Context, cl client.Client, logger logr.Logger, existing *integrationv1beta2.NudgeConfig, newRelationships []integrationv1beta2.NudgeRelationship, dryRun bool, summary migrationSummary) migrationSummary {
	merged, addedCount := mergeNudges(existing.Spec.Nudges, newRelationships)
	summary.RelationshipsMigrated = addedCount
	summary.RelationshipsAlreadyPresent = len(newRelationships) - addedCount

	if addedCount == 0 {
		if dryRun {
			summary.Action = "dry-run:skip"
		} else {
			summary.Action = "skipped"
		}
		logger.Info("All relationships already present in existing NudgeConfig, no changes needed",
			"existingEntries", len(existing.Spec.Nudges))
		return summary
	}

	if len(merged) > MaxNudgesPerConfig {
		err := fmt.Errorf("merged count %d exceeds maximum %d", len(merged), MaxNudgesPerConfig)
		logger.Error(err, "Merged NudgeConfig would exceed max entries, skipping namespace")
		summary.Action = "error"
		summary.Errors = append(summary.Errors, err.Error())
		return summary
	}

	if err := dag.ValidateNudgeGraph(merged); err != nil {
		logger.Error(err, "Cycle detected in merged nudge relationships, skipping namespace")
		summary.Action = "error"
		summary.Errors = append(summary.Errors, fmt.Sprintf("cycle in merged graph: %v", err))
		return summary
	}

	if dryRun {
		summary.Action = "dry-run:update"
		logger.Info("Would UPDATE NudgeConfig",
			"existingEntries", len(existing.Spec.Nudges),
			"entriesToAdd", addedCount)
		return summary
	}

	existing.Spec.Nudges = merged
	setMigrationMetadata(existing)
	if err := cl.Update(ctx, existing); err != nil {
		logger.Error(err, "Failed to update NudgeConfig")
		summary.Action = "error"
		summary.Errors = append(summary.Errors, fmt.Sprintf("update nudgeconfig: %v", err))
		return summary
	}

	summary.Action = "updated"
	logger.Info("Updated NudgeConfig",
		"totalEntries", len(merged),
		"entriesAdded", addedCount)
	return summary
}

func buildNudgeConfig(namespace string, nudges []integrationv1beta2.NudgeRelationship) *integrationv1beta2.NudgeConfig {
	nc := &integrationv1beta2.NudgeConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      integrationv1beta2.NudgeConfigSingletonName,
			Namespace: namespace,
		},
		Spec: integrationv1beta2.NudgeConfigSpec{
			Nudges: nudges,
		},
	}
	setMigrationMetadata(nc)
	return nc
}

func setMigrationMetadata(nc *integrationv1beta2.NudgeConfig) {
	if nc.Labels == nil {
		nc.Labels = make(map[string]string)
	}
	nc.Labels[MigrationLabel] = MigrationLabelValue

	if nc.Annotations == nil {
		nc.Annotations = make(map[string]string)
	}
	nc.Annotations[MigrationAnnotation] = MigrationAnnotationValue
}

func mergeNudges(existing, incoming []integrationv1beta2.NudgeRelationship) ([]integrationv1beta2.NudgeRelationship, int) {
	existingSet := make(map[string]bool, len(existing))
	for _, r := range existing {
		existingSet[r.From+"->"+r.To] = true
	}

	merged := make([]integrationv1beta2.NudgeRelationship, len(existing))
	copy(merged, existing)

	added := 0
	for _, r := range incoming {
		key := r.From + "->" + r.To
		if !existingSet[key] {
			merged = append(merged, r)
			existingSet[key] = true
			added++
		}
	}

	return merged, added
}

func printSummary(summaries []migrationSummary, logger logr.Logger, dryRun bool) {
	created, updated, skipped, errored := 0, 0, 0, 0
	totalMigrated := 0

	for _, s := range summaries {
		switch s.Action {
		case "created", "dry-run:create":
			created++
		case "updated", "dry-run:update":
			updated++
		case "skipped", "dry-run:skip":
			skipped++
		case "error":
			errored++
		}
		totalMigrated += s.RelationshipsMigrated

		if s.Action == "error" {
			for _, e := range s.Errors {
				logger.Error(nil, "Namespace error", "namespace", s.Namespace, "error", e)
			}
		}
	}

	prefix := ""
	if dryRun {
		prefix = "[DRY RUN] "
	}

	logger.Info(fmt.Sprintf("%sMigration summary", prefix),
		"namespacesProcessed", len(summaries),
		"created", created,
		"updated", updated,
		"skipped", skipped,
		"errors", errored,
		"totalRelationshipsMigrated", totalMigrated,
	)
}
