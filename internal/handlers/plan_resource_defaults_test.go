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

func TestResourceDefaultsError(t *testing.T) {
	cases := []struct {
		name string
		plan models.Plan
		ok   bool
	}{
		{"no defaults", models.Plan{MaxMemoryMB: 0, MaxCPUCores: 0}, true},
		{"fits under the totals", models.Plan{DefaultMemoryMB: 256, DefaultCPUMillicores: 500, MaxMemoryMB: 4096, MaxCPUCores: 4}, true},
		{"unlimited totals", models.Plan{DefaultMemoryMB: 8192, DefaultCPUMillicores: 8000, MaxMemoryMB: -1, MaxCPUCores: -1}, true},
		{"memory below the floor", models.Plan{DefaultMemoryMB: 8, MaxMemoryMB: -1}, false},
		{"memory over the total", models.Plan{DefaultMemoryMB: 2048, MaxMemoryMB: 1024}, false},
		{"cpu over the total", models.Plan{DefaultCPUMillicores: 3000, MaxCPUCores: 2}, false},
		{"memory default on a plan with no memory", models.Plan{DefaultMemoryMB: 64, MaxMemoryMB: 0}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := resourceDefaultsError(&tc.plan); (err == nil) != tc.ok {
				t.Errorf("err = %v, want ok=%v", err, tc.ok)
			}
		})
	}
}
