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

package wsbackup

import (
	"testing"

	"github.com/miabi-io/miabi/internal/models"
	"github.com/miabi-io/miabi/internal/wsbundle"
)

// A backup taken before "none" disabled the image's own check must restore to what it ran with.
func TestRestoredHealthcheckType(t *testing.T) {
	cases := []struct {
		typ    string
		schema int
		want   models.HealthcheckType
	}{
		{"none", 2, models.HealthcheckImage},
		{"none", wsbundle.HealthcheckNoneDisables, models.HealthcheckNone},
		{"http", 2, models.HealthcheckHTTP},
		{"image", wsbundle.StateSchema, models.HealthcheckImage},
	}
	for _, c := range cases {
		if got := restoredHealthcheckType(c.typ, c.schema); got != c.want {
			t.Errorf("%q at schema %d = %q, want %q", c.typ, c.schema, got, c.want)
		}
	}
}
