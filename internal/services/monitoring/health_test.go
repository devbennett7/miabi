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

package monitoring

import (
	"testing"

	"github.com/miabi-io/miabi/internal/models"
)

func TestHealthOf(t *testing.T) {
	cases := []struct {
		name       string
		status     models.AppStatus
		hasCurrent bool
		want       string
	}{
		{"running is healthy", models.AppStatusRunning, true, "healthy"},
		{"failed is unhealthy", models.AppStatusFailed, true, "unhealthy"},
		{"deploying over a live release stays healthy", models.AppStatusDeploying, true, "healthy"},
		{"first-ever deploy is unknown (nothing running)", models.AppStatusDeploying, false, "unknown"},
		{"stopped is unknown", models.AppStatusStopped, true, "unknown"},
		{"created is unknown", models.AppStatusCreated, false, "unknown"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := healthOf(tc.status, tc.hasCurrent); got != tc.want {
				t.Errorf("healthOf(%q, %v) = %q, want %q", tc.status, tc.hasCurrent, got, tc.want)
			}
		})
	}
}
