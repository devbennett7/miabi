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
	"strings"
	"testing"

	"github.com/miabi-io/miabi/internal/declarative"
	"github.com/miabi-io/miabi/internal/models"
)

func TestGitAppProjectsItsSourceNotItsBuiltImage(t *testing.T) {
	app := &models.Application{
		Name: "web", SourceType: models.AppSourceGit,
		Image: "registry.example.com/ws1/web", Tag: "sha-9f2c1ab",
		GitRepo: "https://github.com/acme/web", GitRef: "main",
		BuildMethod: models.BuildBuildpack, Builder: "paketobuildpacks/builder-jammy-base",
	}
	r := appResource(app, nil, nil, nil, nil, nil)
	spec := r.Application
	if spec.Image != "" || spec.Tag != "" {
		t.Errorf("the built image leaked into the manifest: %q:%q", spec.Image, spec.Tag)
	}
	if spec.Source == nil || spec.Source.Git != "https://github.com/acme/web" || spec.Source.Ref != "main" {
		t.Fatalf("source = %+v", spec.Source)
	}
	if spec.Source.BuildMethod != "buildpack" || spec.Source.Builder == "" {
		t.Errorf("build config lost: %+v", spec.Source)
	}
	b, err := declarative.Marshal(setOf(r))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "source:") || strings.Contains(string(b), "sha-9f2c1ab") {
		t.Errorf("rendered bundle is wrong:\n%s", b)
	}
}

func TestImageAppKeepsItsImage(t *testing.T) {
	app := &models.Application{Name: "web", SourceType: models.AppSourceImage, Image: "nginx", Tag: "1.27"}
	spec := appResource(app, nil, nil, nil, nil, nil).Application
	if spec.Source != nil {
		t.Errorf("an image app grew a source block: %+v", spec.Source)
	}
	if spec.Image != "nginx" || spec.Tag != "1.27" {
		t.Errorf("image = %q:%q", spec.Image, spec.Tag)
	}
}

func setOf(r declarative.Resource) *declarative.ResourceSet {
	s := declarative.NewResourceSet()
	s.Add(r)
	return s
}
