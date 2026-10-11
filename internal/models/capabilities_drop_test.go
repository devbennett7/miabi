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
	"reflect"
	"testing"
)

// Any Linux capability may be dropped, not only the grantable ones: dropping takes privilege away.
func TestNormalizeDropCapabilities(t *testing.T) {
	got, err := NormalizeDropCapabilities([]string{" net_raw", "CAP_SYS_CHROOT", "NET_RAW", "SYS_BOOT"})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"NET_RAW", "SYS_BOOT", "SYS_CHROOT"}; !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, _ := NormalizeDropCapabilities([]string{"NET_RAW", "all"}); !reflect.DeepEqual(got, []string{CapabilityAll}) {
		t.Errorf("got %v, want ALL alone", got)
	}
	if _, err := NormalizeDropCapabilities([]string{"NET_RAWW"}); !errors.Is(err, ErrCapabilityNotLinux) {
		t.Errorf("a typo was accepted (err = %v)", err)
	}
}

func TestCheckCapabilityConflict(t *testing.T) {
	if err := CheckCapabilityConflict([]string{"NET_ADMIN"}, []string{"NET_ADMIN"}); !errors.Is(err, ErrCapabilityConflict) {
		t.Errorf("err = %v, want a conflict", err)
	}
	if err := CheckCapabilityConflict([]string{"NET_ADMIN"}, []string{CapabilityAll}); err != nil {
		t.Errorf("dropping ALL to keep only the grants was refused: %v", err)
	}
}
