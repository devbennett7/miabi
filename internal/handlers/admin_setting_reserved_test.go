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

package handlers

import "testing"

func TestReservedSettingsAreHiddenFromTheGenericList(t *testing.T) {
	reserved := []string{
		"app.install_id",
		"image.postgres",
		"brand.name",
		"brand.links",
		"brand.accent",
		"brand.logo_url",
		"brand.logo_dark_url",
		"brand.accent_policy",
		"brand.signin_notice",
		"cluster_name",
	}
	for _, key := range reserved {
		if !isReservedSetting(key) {
			t.Errorf("%q is editable in the generic settings list, but has its own screen", key)
		}
	}

	editable := []string{
		"maintenance_mode",
		"default_workspace_role",
		"custom_labels_enabled",
		"repo_pipelines_enabled",
		"require_email_verification",
		"allowed_signup_domains",
		"max_memory_mb",
		"branding_notes",
		"cluster_names",
		"imagemagick_enabled",
	}
	for _, key := range editable {
		if isReservedSetting(key) {
			t.Errorf("%q was reserved, but it belongs in the settings list", key)
		}
	}
}
