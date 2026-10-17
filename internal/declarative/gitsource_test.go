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
