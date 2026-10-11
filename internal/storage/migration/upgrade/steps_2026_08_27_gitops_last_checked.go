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

func gitopsLastCheckedStep(ctx context.Context, db *gorm.DB) error {
	if !db.Migrator().HasColumn(&gitSourceRow{}, "last_checked_at") {
		return nil
	}
	err := db.WithContext(ctx).Exec(`UPDATE git_sources
		SET last_checked_at = last_synced_at, last_checked_commit = last_synced_commit
		WHERE last_checked_at IS NULL AND last_synced_at IS NOT NULL`).Error
	if err != nil {
		return fmt.Errorf("seed git_sources.last_checked_at: %w", err)
	}
	return nil
}

// gitSourceRow names the table for the migrator without importing models.
type gitSourceRow struct{}

func (gitSourceRow) TableName() string { return "git_sources" }

func init() {
	steps = append(steps, Step{
		Name:    "gitops_last_checked_backfill",
		Version: "1.9.0",
		Run:     gitopsLastCheckedStep,
	})
}
