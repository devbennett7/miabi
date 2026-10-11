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

package nodestats

import "testing"

func TestSampleConsistency(t *testing.T) {
	const gb = int64(1) << 30

	cases := []struct {
		name    string
		sampled uint64
		docker  int64
		want    bool
	}{
		// An ordinary host: docker info MemTotal IS /proc's MemTotal, to the byte.
		{"bare metal node", uint64(12599369728), 12599369728, true},
		{"rounding drift is tolerated", uint64(12599369728 - 100), 12599369728, true},
		// The case that broke the dashboard: a 4 GB node on a 64 GB machine reports the machine.
		{"containerised node reports its host", uint64(64 * gb), 4 * gb, false},
		{"memory-limited vm", uint64(12 * gb), 512 * 1024 * 1024, false},
		// No figure from Docker: nothing to contradict the sample, so keep it.
		{"docker reports nothing", uint64(8 * gb), 0, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := describesNode(tc.sampled, tc.docker); got != tc.want {
				t.Fatalf("describesNode(%d, %d) = %v, want %v", tc.sampled, tc.docker, got, tc.want)
			}
		})
	}
}
