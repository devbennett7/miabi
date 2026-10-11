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

import "testing"

func TestDeriveLiveStatus(t *testing.T) {
	cases := []struct {
		name       string
		state      string
		health     string
		restarting bool
		stopped    bool
		want       string
	}{
		{"running healthy", "running", "healthy", false, false, "running"},
		{"running no healthcheck", "running", "", false, false, "running"},
		{"running unhealthy", "running", "unhealthy", false, false, "unhealthy"},
		{"running starting", "running", "starting", false, false, "starting"},
		{"running but restarting flag", "running", "", true, false, "restarting"},
		{"restarting state", "restarting", "", false, false, "restarting"},
		{"exited unexpectedly", "exited", "", false, false, "exited"},
		{"exited after stop", "exited", "", false, true, "stopped"},
		{"dead crash", "dead", "", false, false, "exited"},
		{"paused", "paused", "", false, false, "paused"},
		{"created", "created", "", false, false, "created"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := deriveLiveStatus(c.state, c.health, c.restarting, c.stopped); got != c.want {
				t.Errorf("deriveLiveStatus(%q,%q,%v,%v) = %q, want %q", c.state, c.health, c.restarting, c.stopped, got, c.want)
			}
		})
	}
}
