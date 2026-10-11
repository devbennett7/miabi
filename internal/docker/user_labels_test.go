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

package docker

import "testing"

func TestIsReservedLabelKey(t *testing.T) {
	reserved := []string{
		"io.miabi.app",
		"io.miabi.workspace",
		"com.docker.compose.project",
		"com.docker.anything",
	}
	for _, k := range reserved {
		if !IsReservedLabelKey(k) {
			t.Errorf("expected %q to be reserved", k)
		}
	}
	allowed := []string{
		"traefik.enable",
		"traefik.http.routers.blog.rule",
		"org.opencontainers.image.source",
		"com.example.team",
		"autoheal",
	}
	for _, k := range allowed {
		if IsReservedLabelKey(k) {
			t.Errorf("expected %q to be allowed", k)
		}
	}
}

func TestSanitizeUserLabels(t *testing.T) {
	if got := SanitizeUserLabels(nil); got != nil {
		t.Fatalf("nil in should give nil out, got %v", got)
	}

	in := map[string]string{
		"traefik.enable":              "true",
		"traefik.http.routers.x.rule": "Host(`x`)",
		"io.miabi.app":                "999",       // spoof attempt — must be stripped
		"io.miabi.workspace":          "7",         // spoof attempt — must be stripped
		"io.miabi.managed":            "true",      // spoof attempt — must be stripped
		"com.docker.compose.project":  "evil",      // compose key — stripped
		"":                            "empty-key", // empty key — dropped
	}
	out := SanitizeUserLabels(in)

	if out["traefik.enable"] != "true" || out["traefik.http.routers.x.rule"] != "Host(`x`)" {
		t.Errorf("user labels should survive, got %v", out)
	}
	for _, k := range []string{"io.miabi.app", "io.miabi.workspace", "io.miabi.managed", "com.docker.compose.project", ""} {
		if _, ok := out[k]; ok {
			t.Errorf("reserved/empty key %q must be stripped, got %v", k, out)
		}
	}
	if len(out) != 2 {
		t.Errorf("expected 2 surviving labels, got %d: %v", len(out), out)
	}

	// The input map must not be mutated.
	if _, ok := in["io.miabi.app"]; !ok {
		t.Error("SanitizeUserLabels must not mutate its input")
	}
}
