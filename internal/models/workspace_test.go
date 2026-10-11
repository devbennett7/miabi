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

package models

import "testing"

func TestWorkspaceRoleAtLeast(t *testing.T) {
	cases := []struct {
		role WorkspaceRole
		min  WorkspaceRole
		want bool
	}{
		{WorkspaceRoleOwner, WorkspaceRoleAdmin, true},
		{WorkspaceRoleAdmin, WorkspaceRoleAdmin, true},
		{WorkspaceRoleDeveloper, WorkspaceRoleAdmin, false},
		{WorkspaceRoleViewer, WorkspaceRoleDeveloper, false},
		{WorkspaceRoleDeveloper, WorkspaceRoleViewer, true},
		{WorkspaceRole("bogus"), WorkspaceRoleViewer, false},
	}
	for _, c := range cases {
		if got := c.role.AtLeast(c.min); got != c.want {
			t.Errorf("%q.AtLeast(%q) = %v, want %v", c.role, c.min, got, c.want)
		}
	}
}

func TestWorkspaceRoleValid(t *testing.T) {
	if !WorkspaceRoleOwner.Valid() {
		t.Error("owner should be valid")
	}
	if WorkspaceRole("manager").Valid() {
		t.Error("manager should be invalid")
	}
}
