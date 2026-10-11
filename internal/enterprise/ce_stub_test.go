//go:build !enterprise

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
	"errors"
	"testing"
)

// TestCommunityLDAPDisabled verifies the Community build links no LDAP
// authenticator (so the go-ldap client is never compiled in) and that the
// sso_ldap flag is denied.
func TestCommunityLDAPDisabled(t *testing.T) {
	ee := New(nil, "", "", "")
	if ee.LDAP() != nil {
		t.Error("Community LDAP() must be nil")
	}
	if err := ee.Require(FlagSSOLDAP); !errors.Is(err, ErrLicenseRequired) {
		t.Errorf("Require(sso_ldap) = %v, want ErrLicenseRequired", err)
	}
	if ee.Has(FlagSSOLDAP) {
		t.Error("Community must not have sso_ldap")
	}
}

// TestSSOLDAPFlagRegistered ensures the flag is in the canonical AllFlags list
// (so the issuer tool + docs advertise it).
func TestSSOLDAPFlagRegistered(t *testing.T) {
	if !IsKnownFlag(FlagSSOLDAP) {
		t.Errorf("%q missing from AllFlags", FlagSSOLDAP)
	}
}
