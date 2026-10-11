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

package middlewares

import "testing"

type fakeApps map[string]uint

func (f fakeApps) IDByUID(uid string) (uint, error) {
	if id, ok := f[uid]; ok {
		return id, nil
	}
	return 0, errNoApp
}

var errNoApp = errAppNotFound{}

type errAppNotFound struct{}

func (errAppNotFound) Error() string { return "not found" }

func TestResolveAppRef(t *testing.T) {
	apps := fakeApps{"01J8ZK": 7}
	cases := []struct {
		ref  string
		want uint
		ok   bool
	}{
		{"7", 7, true},
		{"01J8ZK", 7, true},
		{"01OTHER", 0, false}, // an unknown uid never resolves, so the caller refuses
		{"0", 0, false},
		{"-1", 0, false},
		{"", 0, false},
		{"not-a-uid", 0, false},
	}
	for _, c := range cases {
		got, ok := resolveAppRef(c.ref, apps)
		if got != c.want || ok != c.ok {
			t.Errorf("resolveAppRef(%q) = (%d, %v), want (%d, %v)", c.ref, got, ok, c.want, c.ok)
		}
	}
}

// Without a resolver an opaque reference cannot be checked, so it must not pass.
func TestResolveAppRefWithoutResolver(t *testing.T) {
	if _, ok := resolveAppRef("01J8ZK", nil); ok {
		t.Fatal("an unresolvable reference must not be accepted")
	}
	// A numeric id needs no resolver.
	if id, ok := resolveAppRef("7", nil); !ok || id != 7 {
		t.Fatalf("numeric ref = (%d, %v), want (7, true)", id, ok)
	}
}
