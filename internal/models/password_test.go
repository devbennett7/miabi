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
package models

import (
	"errors"
	"strings"
	"testing"
)

func TestCheckNewPasswordCountsCharacters(t *testing.T) {
	for pw, ok := range map[string]bool{
		strings.Repeat("a", 11): false,
		strings.Repeat("a", 12): true,
		strings.Repeat("é", 12): true, // 24 bytes, 12 characters
		strings.Repeat("é", 11): false,
	} {
		err := CheckNewPassword(pw)
		if ok != (err == nil) || (err != nil && !errors.Is(err, ErrPasswordTooShort)) {
			t.Errorf("CheckNewPassword(%d runes) = %v, want ok=%v", len([]rune(pw)), err, ok)
		}
	}
}
