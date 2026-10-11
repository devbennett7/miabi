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

package manifest

import (
	"strings"
	"testing"
)

const sizedDatabaseYAML = `
apiVersion: miabi.io/v1
kind: Template
metadata:
  name: shop
  displayName: Shop
  version: 1.0.0
databases:
  - name: db
    engine: postgres
    placement: %s
    resources:
      memory: %s
applications:
  - name: shop
    image: nginx
`

func sizedDatabase(placement, memory string) []byte {
	return []byte(strings.NewReplacer("placement: %s", "placement: "+placement, "memory: %s", "memory: "+memory).Replace(sizedDatabaseYAML))
}

func TestDatabaseResourcesParse(t *testing.T) {
	m, err := Parse(sizedDatabase("auto", "1Gi"))
	if err != nil {
		t.Fatal(err)
	}
	if r := m.Databases[0].Resources; r == nil || r.Memory != "1Gi" {
		t.Errorf("resources = %+v", r)
	}
}

func TestDatabaseResourcesValidate(t *testing.T) {
	for _, tc := range []struct{ name, placement, memory, want string }{
		{"bad memory", "auto", "plenty", "invalid memory"},
		{"shared", "shared", "1Gi", "instance of its own"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := Parse(sizedDatabase(tc.placement, tc.memory)); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("err = %v, want it to mention %q", err, tc.want)
			}
		})
	}
}
