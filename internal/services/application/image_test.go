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

func TestNormalizeImageTag(t *testing.T) {
	cases := []struct {
		name      string
		image     string
		tag       string
		wantImage string
		wantTag   string
	}{
		{"image with tag, no tag field", "nginx:1.2", "", "nginx", "1.2"},
		{"bare image, tag field", "nginx", "v2", "nginx", "v2"},
		{"bare image, no tag", "nginx", "", "nginx", ""},
		{"registry path with tag", "ghcr.io/org/app:v2", "", "ghcr.io/org/app", "v2"},
		{"host:port and tag", "registry:5000/app:1.2.3", "", "registry:5000/app", "1.2.3"},
		{"embedded tag wins over tag field", "nginx:1.2", "v9", "nginx", "1.2"},
		{"digest pin kept intact", "repo@sha256:abc", "", "repo@sha256:abc", ""},
		{"whitespace trimmed", "  nginx:1.2  ", "  ", "nginx", "1.2"},
		{"empty image", "", "v2", "", "v2"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			img, tag := normalizeImageTag(c.image, c.tag)
			if img != c.wantImage || tag != c.wantTag {
				t.Errorf("normalizeImageTag(%q, %q) = (%q, %q), want (%q, %q)", c.image, c.tag, img, tag, c.wantImage, c.wantTag)
			}
		})
	}
}
