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

// Plan placement is sold from the Business tier up, so the flag must be grantable and those tiers must grant it.
func TestPlacementPolicyEntitlement(t *testing.T) {
	if !IsKnownFlag(FlagPlacementPolicy) {
		t.Fatalf("%s is missing from AllFlags, so no license can grant it", FlagPlacementPolicy)
	}
	for _, tier := range []string{TierBusiness, TierEnterprise} {
		p, ok := TierByName(tier)
		if !ok {
			t.Fatalf("tier %q is not defined", tier)
		}
		if !slices.Contains(p.Flags, FlagPlacementPolicy) {
			t.Errorf("tier %q does not grant %s", tier, FlagPlacementPolicy)
		}
	}
	if p, _ := TierByName(TierProfessional); slices.Contains(p.Flags, FlagPlacementPolicy) {
		t.Errorf("tier %q grants %s", TierProfessional, FlagPlacementPolicy)
	}
}
