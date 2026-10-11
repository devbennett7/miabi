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

package worker

import (
	"reflect"
	"testing"

	"github.com/miabi-io/miabi/internal/models"
)

// An app's hardening can only add to the profile's: neither side's restriction may be lost.
func TestAppHardeningTightensTheProfile(t *testing.T) {
	profile := Security{NoNewPrivileges: true, CapDrop: []string{"NET_RAW"}, Restricted: true}
	got := profile.withHardening(&models.Application{ReadOnlyRootFilesystem: true, DropCapabilities: []string{"NET_RAW", "SYS_CHROOT"}})
	if !got.NoNewPrivileges || !got.ReadOnlyRootfs || !reflect.DeepEqual(got.CapDrop, []string{"NET_RAW", "SYS_CHROOT"}) {
		t.Errorf("security = %+v", got)
	}
	if len(profile.CapDrop) != 1 {
		t.Errorf("the profile's drop list was modified: %v", profile.CapDrop)
	}
	if got := profile.withHardening(&models.Application{}); !got.NoNewPrivileges || got.ReadOnlyRootfs || len(got.CapDrop) != 1 {
		t.Errorf("an app asking for nothing loosened the profile: %+v", got)
	}
	if got := (Security{}).withHardening(&models.Application{NoNewPrivileges: true}); !got.NoNewPrivileges {
		t.Error("an app's no-new-privileges was dropped under the default profile")
	}
}
