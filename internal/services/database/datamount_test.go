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

package database

import (
	"testing"

	"github.com/miabi-io/miabi/internal/models"
)

func TestPostgresMajor(t *testing.T) {
	cases := map[string]int{
		"18":          18,
		"18.1":        18,
		"18-alpine":   18,
		"17-bookworm": 17,
		"17-alpine":   17,
		"9.6":         9,
		"":            0,
		"latest":      0,
		"alpine":      0,
		"180":         180,
		"18beta1":     18,
	}
	for in, want := range cases {
		if got := postgresMajor(in); got != want {
			t.Errorf("postgresMajor(%q) = %d, want %d", in, got, want)
		}
	}
}

func TestDataMount(t *testing.T) {
	const pgDefault = "/var/lib/postgresql/data"

	cases := []struct {
		name    string
		engine  models.DBEngine
		version string
		def     string
		want    string
	}{
		{"postgres 18 mounts one level up", models.DBEnginePostgres, "18-alpine", pgDefault, "/var/lib/postgresql"},
		{"postgres 19 (future) also new layout", models.DBEnginePostgres, "19", pgDefault, "/var/lib/postgresql"},
		{"postgres 17 keeps legacy path", models.DBEnginePostgres, "17-alpine", pgDefault, pgDefault},
		{"postgres unparseable version stays legacy", models.DBEnginePostgres, "latest", pgDefault, pgDefault},
		{"mysql unaffected", models.DBEngineMySQL, "8.4", "/var/lib/mysql", "/var/lib/mysql"},
		{"redis unaffected", models.DBEngineRedis, "7-alpine", "/data", "/data"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := dataMount(c.engine, c.version, c.def); got != c.want {
				t.Errorf("dataMount(%s, %q) = %q, want %q", c.engine, c.version, got, c.want)
			}
		})
	}
}
