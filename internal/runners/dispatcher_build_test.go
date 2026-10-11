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

package runners

import (
	"strings"
	"testing"
	"time"

	"github.com/miabi-io/runner/proto"
)

func TestBuildOnlyJobSpec(t *testing.T) {
	in := BuildInputs{
		DeploymentID:     42,
		DeploymentNumber: 3,
		WorkspaceID:      7,
		WorkspaceName:    "acme",
		AppID:            5,
		AppName:          "web",
		SourceURL:        "https://git.example.com/acme/web.git",
		Commit:           "abc123",
		Repository:       "reg.example.com/ws_7/app-5",
		Registry:         "reg.example.com",
		Build:            &proto.BuildConfig{Method: "buildpack", Builder: "paketobuildpacks/builder-jammy-base"},
	}
	spec, secrets := buildOnlyJobSpec(in, nil, time.Time{})

	// RunID carries the per-app deployment number (the image tag), not the id.
	if spec.RunID != 3 {
		t.Errorf("RunID = %d, want the deployment number 3", spec.RunID)
	}
	if len(spec.Steps) != 1 || spec.Steps[0].Uses != "build" {
		t.Fatalf("want exactly one build step, got %+v", spec.Steps)
	}
	if spec.Steps[0].Build == nil || spec.Steps[0].Build.Method != "buildpack" {
		t.Errorf("build config not carried onto the step: %+v", spec.Steps[0].Build)
	}
	env := strings.Join(spec.Env, "\n")
	for _, want := range []string{
		"MIABI_IMAGE_REPOSITORY=reg.example.com/ws_7/app-5",
		"MIABI_APP_ID=5",
		"MIABI_COMMIT=abc123",
		"MIABI_REGISTRY=reg.example.com",
	} {
		if !strings.Contains(env, want) {
			t.Errorf("env missing %q; got %v", want, spec.Env)
		}
	}
	if len(secrets) != 0 {
		t.Errorf("no creds -> no secrets to mask, got %v", secrets)
	}
}
