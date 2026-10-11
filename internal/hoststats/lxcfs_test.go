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

package hoststats

import (
	"strings"
	"testing"
)

// Taken from the control plane container on an LXC test VM, where /proc is bound to /host/proc.
const lxcMountinfo = `5691 5663 0:112 / /host/proc ro,nosuid,nodev,noexec,relatime - proc proc rw
5697 5691 0:44 /proc/cpuinfo /host/proc/cpuinfo ro,nosuid,nodev,relatime - fuse.lxcfs lxcfs rw,user_id=0,group_id=0,allow_other
5700 5691 0:44 /proc/meminfo /host/proc/meminfo ro,nosuid,nodev,relatime - fuse.lxcfs lxcfs rw,user_id=0,group_id=0,allow_other
5703 5691 0:44 /proc/stat /host/proc/stat ro,nosuid,nodev,relatime - fuse.lxcfs lxcfs rw,user_id=0,group_id=0,allow_other
`

const plainMountinfo = `5691 5663 0:112 / /host/proc ro,nosuid,nodev,noexec,relatime - proc proc rw
812 5691 0:44 / /proc rw,nosuid,nodev,noexec,relatime shared:12 - proc proc rw
`

func TestLXCFSDetection(t *testing.T) {
	cases := []struct {
		name, info, proc string
		want             bool
	}{
		{"lxc guest, bound procfs", lxcMountinfo, "/host/proc", true},
		{"lxc guest, trailing slash", lxcMountinfo, "/host/proc/", true},
		{"lxc guest, other path", lxcMountinfo, "/proc", false},
		{"ordinary host", plainMountinfo, "/host/proc", false},
		{"garbage", "not mountinfo\n", "/host/proc", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := lxcfsMounted(strings.NewReader(tc.info), tc.proc); got != tc.want {
				t.Fatalf("lxcfsMounted = %v, want %v", got, tc.want)
			}
		})
	}
}
