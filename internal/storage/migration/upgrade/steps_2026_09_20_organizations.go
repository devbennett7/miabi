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

	"github.com/jkaninda/logger"
	"gorm.io/gorm"
)

// organizationsStep attaches existing workspaces and users to the default organization.
//
// "NULL means the default org" has been the documented convention since the realm was introduced, and
// reads stay tolerant of it. Populating the column matters because organizations now carry a
// workspace cap and may own clusters: both are counted with a plain `WHERE organization_id = ?`, and
// a fleet of NULLs would make an org that holds every workspace look empty.
//
// Idempotent: it only touches rows that are still NULL.
func organizationsStep(ctx context.Context, db *gorm.DB) error {
	var defaultID uint
	if err := db.WithContext(ctx).Table("organizations").
		Select("id").Where("is_default = ?", true).
		Limit(1).Scan(&defaultID).Error; err != nil {
		return err
	}
	if defaultID == 0 {
		// Seeding runs before upgrade steps, so this means the table is not there yet — a fresh
		// install, which has nothing to backfill.
		logger.Info("organization backfill: no default organization yet, nothing to attach")
		return nil
	}
	for _, table := range []string{"workspaces", "users"} {
		res := db.WithContext(ctx).Table(table).
			Where("organization_id IS NULL").
			Update("organization_id", defaultID)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected > 0 {
			logger.Info("attached rows to the default organization",
				"table", table, "rows", res.RowsAffected, "organization", defaultID)
		}
	}
	return nil
}

func init() {
	steps = append(steps, Step{
		Name:    "organizations_backfill",
		Version: "1.10.4",
		Run:     organizationsStep,
	})
}
