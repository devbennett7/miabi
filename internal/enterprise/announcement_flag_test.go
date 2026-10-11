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

package enterprise

import (
	"slices"
	"testing"
)

// Announcements are a paid capability, so the flag must be grantable and the tier
// presets must actually grant it — otherwise the admin UI stays locked for
// customers whose licence is supposed to include it.
func TestAnnouncementsEntitlement(t *testing.T) {
	if !IsKnownFlag(FlagAnnouncements) {
		t.Fatalf("%s is missing from AllFlags, so no license can grant it", FlagAnnouncements)
	}
	for _, tier := range []string{TierBusiness, TierEnterprise} {
		p, ok := TierByName(tier)
		if !ok {
			t.Fatalf("tier %q is not defined", tier)
		}
		if !slices.Contains(p.Flags, FlagAnnouncements) {
			t.Errorf("tier %q does not grant %s", tier, FlagAnnouncements)
		}
	}
}
