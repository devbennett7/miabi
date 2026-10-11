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

func TestResourcePolicyGrants(t *testing.T) {
	p := &ResourcePolicy{Permissions: []string{string(PermAppRead), string(PermAppDeploy)}}
	if !p.Grants(PermAppDeploy) {
		t.Fatal("expected app:deploy granted")
	}
	if p.Grants(PermAppDelete) {
		t.Fatal("app:delete must not be granted")
	}
	empty := &ResourcePolicy{}
	if empty.Grants(PermAppRead) {
		t.Fatal("empty policy grants nothing")
	}
}

func TestIsValidResourceType(t *testing.T) {
	for _, ok := range []string{ResourceTypeApp, ResourceTypeDatabase, ResourceTypeDomain} {
		if !IsValidResourceType(ok) {
			t.Fatalf("%q should be valid", ok)
		}
	}
	for _, bad := range []string{"", "workspace", "node"} {
		if IsValidResourceType(bad) {
			t.Fatalf("%q should be invalid", bad)
		}
	}
}
