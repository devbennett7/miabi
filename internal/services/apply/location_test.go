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
	"errors"
	"strings"
	"testing"

	"github.com/miabi-io/miabi/internal/declarative"
)

func TestRefuseMove(t *testing.T) {
	moved := declarative.Change{Kind: declarative.KindDatabase, Name: "pg", Fields: []declarative.FieldDiff{
		{Field: "version", From: "16", To: "17"},
		{Field: "placement.location", From: "eu-central", To: "eu-east"},
	}}
	err := refuseMove(moved)
	if !errors.Is(err, ErrInvalidManifest) {
		t.Fatalf("err = %v, want ErrInvalidManifest", err)
	}
	for _, want := range []string{`database "pg"`, "eu-central", "eu-east", "Delete it"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("message is missing %q: %v", want, err)
		}
	}
	stay := declarative.Change{Kind: declarative.KindStack, Name: "shop", Fields: []declarative.FieldDiff{{Field: "description"}}}
	if err := refuseMove(stay); err != nil {
		t.Errorf("an update that keeps its location was refused: %v", err)
	}
}

// References are read from the raw env: once rendered, the templates are gone.
func TestRefsAreReadBeforeRendering(t *testing.T) {
	refs := refsOf(&declarative.ApplicationSpec{Env: map[string]string{
		"API":   "{{ .applications.api.url }}",
		"DB":    "postgres://{{ .databases.pg.user }}@{{ .databases.pg.host }}",
		"PLAIN": "http://api:8080",
	}})
	if !refs.apps["api"] || !refs.dbs["pg"] || len(refs.apps) != 1 || len(refs.dbs) != 1 {
		t.Errorf("refs = %+v", refs)
	}
	if got := refsOf(nil); len(got.apps)+len(got.dbs) != 0 {
		t.Errorf("nil spec refs = %+v", got)
	}
}
