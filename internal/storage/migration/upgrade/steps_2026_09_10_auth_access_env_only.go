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

package upgrade

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

var authAccessEnvOnlyKeys = []string{"registration_enabled", "password_reset_enabled"}

func authAccessEnvOnlyStep(ctx context.Context, db *gorm.DB) error {
	err := db.WithContext(ctx).
		Exec(`DELETE FROM settings WHERE key IN ?`, authAccessEnvOnlyKeys).Error
	if err != nil {
		return fmt.Errorf("drop env-only auth settings: %w", err)
	}
	return nil
}

func init() {
	steps = append(steps, Step{
		Name:    "auth_access_env_only",
		Version: "1.9.7",
		Run:     authAccessEnvOnlyStep,
	})
}
