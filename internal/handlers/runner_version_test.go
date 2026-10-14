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
package handlers

import (
	"testing"

	"github.com/miabi-io/miabi/internal/models"
)

func TestAnnotateRunnerVersions(t *testing.T) {
	h := &RunnerHandler{releases: fakeAgentReleases{latest: "v0.3.0"}}
	runners := []models.Runner{
		{Name: "behind", Version: "0.2.0"},
		{Name: "current", Version: "0.3.0"},
		{Name: "dev", Version: "dev"},
		{Name: "silent"},
	}
	h.annotate(runners)
	for _, r := range runners {
		want := ""
		if r.Name == "behind" {
			want = "v0.3.0"
		}
		if r.LatestVersion != want {
			t.Errorf("%s: LatestVersion = %q, want %q", r.Name, r.LatestVersion, want)
		}
	}

	one := models.Runner{Version: "0.1.0"}
	h.annotateOne(&one)
	if one.LatestVersion != "v0.3.0" {
		t.Errorf("annotateOne: LatestVersion = %q", one.LatestVersion)
	}
	(&RunnerHandler{}).annotate(runners[:1])
}
