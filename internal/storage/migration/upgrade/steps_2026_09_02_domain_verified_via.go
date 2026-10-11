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

func domainVerifiedViaStep(ctx context.Context, db *gorm.DB) error {
	if !db.Migrator().HasColumn(&domainRow{}, "verified_via") {
		return nil
	}
	err := db.WithContext(ctx).Exec(`UPDATE domains
		SET verified_via = 'dns_provider'
		WHERE verified = true AND dns_provider_id IS NOT NULL`).Error
	if err != nil {
		return fmt.Errorf("backfill domains.verified_via: %w", err)
	}
	return nil
}

// domainRow names the table for the migrator without importing models.
type domainRow struct{}

func (domainRow) TableName() string { return "domains" }

func init() {
	steps = append(steps, Step{
		Name:    "domain_verified_via_backfill",
		Version: "1.9.4",
		Run:     domainVerifiedViaStep,
	})
}
