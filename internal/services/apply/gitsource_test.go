// Copyright (c) 2026 Jonas Kaninda.
//
// This file is part of Miabi.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package apply

import (
	"context"
	"errors"
	"testing"

	"github.com/miabi-io/miabi/internal/declarative"
	"github.com/miabi-io/miabi/internal/models"
)

type fakeGitSources struct {
	sources  []models.GitSource
	declared []GitSourceInput
	ids      []uint
	removed  []uint
}

func (f *fakeGitSources) List(uint) ([]models.GitSource, error) { return f.sources, nil }
func (f *fakeGitSources) RepositoryNames(uint) (map[uint]string, error) {
	return map[uint]string{4: "gitea"}, nil
}
func (f *fakeGitSources) Declare(_, id uint, in GitSourceInput) (*models.GitSource, error) {
	f.ids, f.declared = append(f.ids, id), append(f.declared, in)
	return &models.GitSource{ID: id, Name: in.Name}, nil
}
func (f *fakeGitSources) Remove(_ context.Context, _, id uint) error {
	f.removed = append(f.removed, id)
	return nil
}

// root(1) <- child(2) <- grandchild(3)
func nestedSources() []models.GitSource {
	return []models.GitSource{
		{ID: 1, Name: "root"},
		{ID: 2, Name: "child", Metadata: models.Metadata{models.MetaManagedBy: ManagedByGitOps, models.MetaGitOpsSource: "1"}},
		{ID: 3, Name: "grandchild", Metadata: models.Metadata{models.MetaManagedBy: ManagedByGitOps, models.MetaGitOpsSource: "2"}},
	}
}

func TestAncestry(t *testing.T) {
	got := ancestry(nestedSources(), "3")
	if len(got) != 3 || got[0] != "3" || got[2] != "1" {
		t.Fatalf("ancestry(3) = %v, want [3 2 1]", got)
	}
	looped := []models.GitSource{
		{ID: 1, Metadata: models.Metadata{models.MetaGitOpsSource: "2"}},
		{ID: 2, Metadata: models.Metadata{models.MetaGitOpsSource: "1"}},
	}
	if got := ancestry(looped, "1"); len(got) != 2 {
		t.Fatalf("a looped chain must stop, got %v", got)
	}
}

func TestCheckNesting(t *testing.T) {
	src := nestedSources()
	if err := checkNesting(src, "", "root", &src[0]); err != nil {
		t.Fatalf("a one-shot apply is never nested: %v", err)
	}
	if err := checkNesting(src, "1", "new", nil); err != nil {
		t.Fatalf("root declaring a new child: %v", err)
	}
	if err := checkNesting(src, "2", "root", &src[0]); !errors.Is(err, ErrInvalidManifest) {
		t.Fatalf("child declaring its parent = %v, want refused", err)
	}
	if err := checkNesting(src, "2", "child", &src[1]); !errors.Is(err, ErrInvalidManifest) {
		t.Fatalf("a project declaring itself = %v, want refused", err)
	}
	if err := checkNesting(src, "3", "deeper", nil); !errors.Is(err, ErrInvalidManifest) {
		t.Fatalf("a fourth level = %v, want refused", err)
	}
}

func TestApplyGitSourceStampsOwner(t *testing.T) {
	f := &fakeGitSources{sources: nestedSources()}
	s := &Service{gitSources: f}
	desired := declarative.Resource{Kind: declarative.KindGitSource, Metadata: declarative.Meta{Name: "shop"},
		GitSource: &declarative.GitSourceSpec{Repo: "https://x/y.git", Ref: "main", Path: ".", SyncPolicy: "auto"}}
	ch := declarative.Change{Action: declarative.ActionCreate, Kind: declarative.KindGitSource, Name: "shop"}

	if err := s.applyGitSource(withOwnerSource(context.Background(), "1"), 1, ch, desired); err != nil {
		t.Fatal(err)
	}
	if m := f.declared[0].Metadata; m[models.MetaManagedBy] != ManagedByGitOps || m[models.MetaGitOpsSource] != "1" {
		t.Fatalf("nested source metadata = %v, want gitops from project 1", m)
	}
	if err := s.applyGitSource(context.Background(), 1, ch, desired); err != nil {
		t.Fatal(err)
	}
	if m := f.declared[1].Metadata; m[models.MetaManagedBy] != ManagedByApply || m[models.MetaGitOpsSource] != "" {
		t.Fatalf("CLI source metadata = %v, want apply with no parent", m)
	}

	del := declarative.Change{Action: declarative.ActionDelete, Kind: declarative.KindGitSource, Name: "child"}
	if err := s.applyGitSource(context.Background(), 1, del, declarative.Resource{}); err != nil || len(f.removed) != 1 || f.removed[0] != 2 {
		t.Fatalf("delete = %v, removed %v, want source 2 removed", err, f.removed)
	}
}

func TestGitSourceSnapshotHasNoDrift(t *testing.T) {
	repoID := uint(4)
	live := models.GitSource{Name: "shop", RepoURL: "https://gitea/x.git", GitRepositoryID: &repoID,
		Ref: "main", Path: "apps", SyncPolicy: models.GitSyncAuto, Prune: true}
	desired, err := declarative.Parse([]byte(`apiVersion: miabi.io/v1
kind: GitSource
metadata: { name: shop }
spec: { repository: gitea, path: apps, syncPolicy: auto, prune: true }
`))
	if err != nil {
		t.Fatal(err)
	}
	actual := declarative.NewResourceSet()
	actual.Add(declarative.Resource{APIVersion: declarative.APIVersion, Kind: declarative.KindGitSource,
		Metadata: declarative.Meta{Name: "shop"}, GitSource: gitSourceSpecOf(&live, map[uint]string{4: "gitea"})})
	plan := declarative.BuildPlan(desired, actual, declarative.PlanOptions{})
	if plan.HasChanges() {
		t.Fatalf("plan = %+v, want no drift", plan.Changes)
	}
}
