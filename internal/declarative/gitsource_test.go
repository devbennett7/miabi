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

package declarative_test

import (
	"strings"
	"testing"

	d "github.com/miabi-io/miabi/internal/declarative"
)

func gitSourceDoc(spec string) string {
	return "apiVersion: miabi.io/v1\nkind: GitSource\nmetadata: { name: shop }\nspec:\n" + spec
}

func TestGitSourceDefaults(t *testing.T) {
	set, err := d.Parse([]byte(gitSourceDoc("  repo: https://git.example.com/shop.git\n")))
	if err != nil {
		t.Fatal(err)
	}
	r, _ := set.Get("GitSource/shop")
	g := r.GitSource
	if g.Ref != "main" || g.Path != "." || g.SyncPolicy != "manual" || g.DeletePolicy != d.DeletePolicyOrphan {
		t.Fatalf("defaults = %+v", g)
	}
}

func TestGitSourceValidation(t *testing.T) {
	cases := map[string]string{
		"no repository":     "  ref: main\n",
		"both repositories": "  repo: https://x/y.git\n  repository: creds\n",
		"bad sync policy":   "  repo: https://x/y.git\n  syncPolicy: always\n",
		"bad delete policy": "  repo: https://x/y.git\n  deletePolicy: purge\n",
		"absolute path":     "  repo: https://x/y.git\n  path: /etc\n",
		"escaping path":     "  repo: https://x/y.git\n  path: apps/../../secrets\n",
	}
	for name, spec := range cases {
		if _, err := d.Parse([]byte(gitSourceDoc(spec))); err == nil || !strings.Contains(err.Error(), "shop") {
			t.Errorf("%s: err = %v, want a gitsource error", name, err)
		}
	}
}

func TestGitSourceRanksAfterItsSiblings(t *testing.T) {
	bundle := gitSourceDoc("  repo: https://x/y.git\n") + `---
apiVersion: miabi.io/v1
kind: Application
metadata: { name: api }
spec: { image: nginx }
`
	desired, err := d.Parse([]byte(bundle))
	if err != nil {
		t.Fatal(err)
	}
	plan := d.BuildPlan(desired, d.NewResourceSet(), d.PlanOptions{})
	if last := plan.Changes[len(plan.Changes)-1]; last.Kind != d.KindGitSource {
		t.Fatalf("last change = %s, want the GitSource", last.Kind)
	}
}

func gitRepositoryDoc(spec string) string {
	return "apiVersion: miabi.io/v1\nkind: GitRepository\nmetadata: { name: gitea }\nspec:\n" + spec
}

func TestGitRepositoryNormalizes(t *testing.T) {
	set, err := d.Parse([]byte(gitRepositoryDoc("  url: https://git.example.com/shop/\n  secretRef: gitea-token\n")))
	if err != nil {
		t.Fatal(err)
	}
	r, _ := set.Get("GitRepository/gitea")
	if g := r.GitRepository; g.URL != "https://git.example.com/shop.git" || g.Auth != "token" {
		t.Fatalf("normalized = %+v, want a .git URL and token auth", g)
	}
	set, err = d.Parse([]byte(gitRepositoryDoc("  url: https://github.com/miabi-io/examples.git\n")))
	if err != nil {
		t.Fatal(err)
	}
	if r, _ := set.Get("GitRepository/gitea"); r.GitRepository.Auth != "public" {
		t.Fatalf("auth without a secretRef = %q, want public", r.GitRepository.Auth)
	}
}

func TestGitRepositoryValidation(t *testing.T) {
	cases := map[string]string{
		"no url":             "  secretRef: tok\n",
		"bad auth":           "  url: https://x/y.git\n  auth: basic\n  secretRef: tok\n",
		"public with secret": "  url: https://x/y.git\n  auth: public\n  secretRef: tok\n",
		"ssh without secret": "  url: git@x:y.git\n  auth: ssh\n",
		"inline token":       "  url: https://x/y.git\n  token: abc\n",
	}
	for name, spec := range cases {
		if _, err := d.Parse([]byte(gitRepositoryDoc(spec))); err == nil {
			t.Errorf("%s: parsed, want an error", name)
		}
	}
}

func TestGitRepositoryComesBeforeItsUsers(t *testing.T) {
	bundle := gitRepositoryDoc("  url: https://x/y.git\n  secretRef: tok\n") + `---
apiVersion: miabi.io/v1
kind: Secret
metadata: { name: tok }
spec: { value: ghp_example }
---
` + gitSourceDoc("  repository: gitea\n")
	desired, err := d.Parse([]byte(bundle))
	if err != nil {
		t.Fatal(err)
	}
	plan := d.BuildPlan(desired, d.NewResourceSet(), d.PlanOptions{})
	pos := map[d.Kind]int{}
	for i, c := range plan.Changes {
		pos[c.Kind] = i
	}
	if !(pos[d.KindSecret] < pos[d.KindGitRepository] && pos[d.KindGitRepository] < pos[d.KindGitSource]) {
		t.Fatalf("order = %+v, want secret, credential, then project", plan.Changes)
	}
	edges := d.Edges(desired)
	want := map[string]bool{"GitSource/shop->GitRepository/gitea": false, "GitRepository/gitea->Secret/tok": false}
	for _, e := range edges {
		if _, ok := want[e.From+"->"+e.To]; ok {
			want[e.From+"->"+e.To] = true
		}
	}
	for k, seen := range want {
		if !seen {
			t.Errorf("missing edge %s", k)
		}
	}
}
