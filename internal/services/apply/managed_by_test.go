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

func TestManagedByFollowsRunSource(t *testing.T) {
	if got := managedBy(context.Background()); got != ManagedByApply {
		t.Fatalf("one-shot apply stamps %q, want %q", got, ManagedByApply)
	}
	if got := managedBy(withOwnerSource(context.Background(), "7")); got != ManagedByGitOps {
		t.Fatalf("GitOps sync stamps %q, want %q", got, ManagedByGitOps)
	}
}

func TestClaim(t *testing.T) {
	cli := context.Background()
	sync := withOwnerSource(context.Background(), "7")
	cases := []struct {
		name       string
		ctx        context.Context
		in         models.Metadata
		owner, src string
	}{
		{"project adopts a CLI-applied app", sync, models.Metadata{models.MetaManagedBy: ManagedByApply}, ManagedByGitOps, "7"},
		{"project labels its unsourced app", sync, models.Metadata{models.MetaManagedBy: ManagedByGitOps}, ManagedByGitOps, "7"},
		{"CLI leaves an unsourced gitops app", cli, models.Metadata{models.MetaManagedBy: ManagedByGitOps}, ManagedByGitOps, ""},
		{"CLI app stays the CLI's", cli, models.Metadata{models.MetaManagedBy: ManagedByApply}, ManagedByApply, ""},
		{"forced CLI leaves a project's app", cli, models.Metadata{models.MetaManagedBy: ManagedByGitOps, models.MetaGitOpsSource: "3"}, ManagedByGitOps, "3"},
		{"another project's app stays", sync, models.Metadata{models.MetaManagedBy: ManagedByGitOps, models.MetaGitOpsSource: "3"}, ManagedByGitOps, "3"},
		{"marketplace app untouched", cli, models.Metadata{models.MetaManagedBy: models.ManagedByMarketplace}, models.ManagedByMarketplace, ""},
	}
	for _, tc := range cases {
		got := claim(tc.ctx, tc.in)
		if got[models.MetaManagedBy] != tc.owner || got[models.MetaGitOpsSource] != tc.src {
			t.Errorf("%s: got (%q, %q), want (%q, %q)", tc.name,
				got[models.MetaManagedBy], got[models.MetaGitOpsSource], tc.owner, tc.src)
		}
	}
}

func ownedApp(name string, labels map[string]string) declarative.Resource {
	return declarative.Resource{
		APIVersion: declarative.APIVersion, Kind: declarative.KindApplication,
		Metadata:    declarative.Meta{Name: name, Labels: labels},
		Application: &declarative.ApplicationSpec{Image: "nginx"},
	}
}

func TestRefuseGitOpsOwned(t *testing.T) {
	actual := declarative.NewResourceSet()
	actual.Add(ownedApp("web", map[string]string{
		declarative.LabelManagedBy: ManagedByGitOps, declarative.LabelGitOpsSource: "3",
	}))
	actual.Add(ownedApp("api", map[string]string{declarative.LabelManagedBy: ManagedByApply}))
	actual.Add(ownedApp("legacy", map[string]string{declarative.LabelManagedBy: ManagedByGitOps}))
	plan := &declarative.Plan{Changes: []declarative.Change{
		{Action: declarative.ActionUpdate, Kind: declarative.KindApplication, Name: "api"},
		{Action: declarative.ActionNoop, Kind: declarative.KindApplication, Name: "web"},
	}}
	if err := refuseGitOpsOwned(plan, actual); err != nil {
		t.Fatalf("no project resource changes, got %v", err)
	}
	for _, name := range []string{"web", "legacy"} {
		plan.Changes[1] = declarative.Change{Action: declarative.ActionUpdate, Kind: declarative.KindApplication, Name: name}
		if err := refuseGitOpsOwned(plan, actual); !errors.Is(err, ErrOwnedByGitOps) {
			t.Fatalf("updating %s = %v, want ErrOwnedByGitOps", name, err)
		}
	}
}

func TestCLIPruneSkipsGitOpsResources(t *testing.T) {
	actual := declarative.NewResourceSet()
	actual.Add(ownedApp("web", map[string]string{
		declarative.LabelManagedBy: ManagedByGitOps, declarative.LabelGitOpsSource: "3",
	}))
	actual.Add(ownedApp("old", map[string]string{declarative.LabelManagedBy: ManagedByApply}))
	plan := declarative.BuildPlan(declarative.NewResourceSet(), actual, declarative.PlanOptions{
		Prune: true, PruneManagedBy: ownerOf(""),
	})
	if len(plan.Changes) != 1 || plan.Changes[0].Name != "old" || plan.Changes[0].Action != declarative.ActionDelete {
		t.Fatalf("CLI prune plan = %+v, want only a delete of old", plan.Changes)
	}
}
