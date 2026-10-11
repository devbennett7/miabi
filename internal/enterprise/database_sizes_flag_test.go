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

// Database sizes are sold from the Business tier up, like plan placement, which they sit beside on a plan.
func TestDatabaseSizesEntitlement(t *testing.T) {
	if !IsKnownFlag(FlagDatabaseSizes) {
		t.Fatalf("%s is missing from AllFlags, so no license can grant it", FlagDatabaseSizes)
	}
	for _, tier := range []string{TierBusiness, TierEnterprise} {
		p, ok := TierByName(tier)
		if !ok {
			t.Fatalf("tier %q is not defined", tier)
		}
		if !slices.Contains(p.Flags, FlagDatabaseSizes) {
			t.Errorf("tier %q does not grant %s", tier, FlagDatabaseSizes)
		}
	}
	if p, _ := TierByName(TierProfessional); slices.Contains(p.Flags, FlagDatabaseSizes) {
		t.Errorf("tier %q grants %s", TierProfessional, FlagDatabaseSizes)
	}
}
