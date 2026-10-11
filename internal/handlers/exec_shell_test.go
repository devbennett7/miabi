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

import "testing"

func TestExecShell(t *testing.T) {
	cases := []struct {
		in, want string
		ok       bool
	}{
		{"", "/bin/sh", true},
		{"/bin/sh", "/bin/sh", true},
		{"/bin/bash", "/bin/bash", true},
		{"bash", "/bin/bash", true},
		{" sh ", "/bin/sh", true},
		{"/usr/bin/python3", "", false},
		{"python3", "", false},
		{"/bin/sh -c id", "", false},
		{"/bin/../usr/bin/python3", "", false},
		{"/bin/bash/", "", false},
	}
	for _, tc := range cases {
		got, ok := execShell(tc.in)
		if got != tc.want || ok != tc.ok {
			t.Errorf("execShell(%q) = %q, %v; want %q, %v", tc.in, got, ok, tc.want, tc.ok)
		}
	}
}
