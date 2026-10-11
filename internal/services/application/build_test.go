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

package application

import (
	"errors"
	"testing"

	"github.com/miabi-io/miabi/internal/models"
)

func TestValidateBuildConfig(t *testing.T) {
	cases := []struct {
		name    string
		source  models.AppSourceType
		method  models.AppBuildMethod
		builder string
		bps     []string
		env     map[string]string
		wantErr error
	}{
		{name: "git auto ok", source: models.AppSourceGit, method: models.BuildAuto},
		{name: "git empty method ok", source: models.AppSourceGit},
		{name: "git buildpack with builder ok", source: models.AppSourceGit, method: models.BuildBuildpack, builder: "x/y"},
		{name: "git invalid method", source: models.AppSourceGit, method: "weird", wantErr: ErrInvalidBuildMethod},

		{name: "image auto allowed (stored default)", source: models.AppSourceImage, method: models.BuildAuto},
		{name: "image empty allowed", source: models.AppSourceImage},
		{name: "image explicit buildpack rejected", source: models.AppSourceImage, method: models.BuildBuildpack, wantErr: ErrBuildConfigOnImage},
		{name: "image with builder rejected", source: models.AppSourceImage, builder: "x/y", wantErr: ErrBuildConfigOnImage},
		{name: "image with buildpacks rejected", source: models.AppSourceImage, bps: []string{"a"}, wantErr: ErrBuildConfigOnImage},
		{name: "image with build env rejected", source: models.AppSourceImage, env: map[string]string{"K": "V"}, wantErr: ErrBuildConfigOnImage},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateBuildConfig(tc.source, tc.method, tc.builder, tc.bps, tc.env)
			if tc.wantErr == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantErr != nil && !errors.Is(err, tc.wantErr) {
				t.Fatalf("err = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestBuildMethodForSource(t *testing.T) {
	if got := buildMethodForSource(models.AppSourceImage, models.BuildBuildpack); got != models.BuildAuto {
		t.Errorf("image source should normalize to auto, got %q", got)
	}
	if got := buildMethodForSource(models.AppSourceGit, ""); got != models.BuildAuto {
		t.Errorf("git empty should default to auto, got %q", got)
	}
	if got := buildMethodForSource(models.AppSourceGit, models.BuildDockerfile); got != models.BuildDockerfile {
		t.Errorf("git explicit should be preserved, got %q", got)
	}
	if got := buildMethodForSource(models.AppSourceGit, "garbage"); got != models.BuildAuto {
		t.Errorf("git invalid should default to auto, got %q", got)
	}
}
