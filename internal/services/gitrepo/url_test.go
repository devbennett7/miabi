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

package gitrepo

import "testing"

func TestNormalizeGitURL(t *testing.T) {
	cases := map[string]string{
		"https://github.com/acme/api":     "https://github.com/acme/api.git",
		"https://github.com/acme/api.git": "https://github.com/acme/api.git",
		"https://github.com/acme/api/":    "https://github.com/acme/api.git",
		"  https://github.com/acme/api  ": "https://github.com/acme/api.git",
		"git@github.com:acme/api":         "git@github.com:acme/api.git",
		"git@github.com:acme/api.git":     "git@github.com:acme/api.git",
		"https://github.com/acme/api.GIT": "https://github.com/acme/api.GIT", // case-insensitive suffix check
		"":                                "",
	}
	for in, want := range cases {
		if got := normalizeGitURL(in); got != want {
			t.Errorf("normalizeGitURL(%q) = %q, want %q", in, got, want)
		}
	}
}
