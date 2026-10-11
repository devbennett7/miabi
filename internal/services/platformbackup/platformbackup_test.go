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

package platformbackup

import "testing"

func TestVolumeArchiveName(t *testing.T) {
	cases := map[string]string{
		"mb-node-gateway-providers": "mb-node-gateway-providers",
		"mb_data":                   "mb_data",
		"weird/name:tag":            "weird-name-tag",
		"":                          "platform-volume",
	}
	for in, want := range cases {
		if got := volumeArchiveName(in); got != want {
			t.Errorf("volumeArchiveName(%q) = %q, want %q", in, got, want)
		}
	}
}
