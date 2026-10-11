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

package healthprobe

import "testing"

func TestGoArch(t *testing.T) {
	for in, want := range map[string]string{"x86_64": "amd64", "amd64": "amd64", "aarch64": "arm64", "ARM64": "arm64"} {
		if got, ok := goArch(in); !ok || got != want {
			t.Errorf("goArch(%q) = %q, %v", in, got, ok)
		}
	}
	if _, ok := goArch("s390x"); ok {
		t.Error("unsupported arch should report false")
	}
	if _, ok := Binary("s390x"); ok {
		t.Error("Binary should be absent for an unsupported arch")
	}
}
