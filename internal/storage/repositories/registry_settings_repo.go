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

package repositories

import (
	"github.com/miabi-io/miabi/internal/models"
	"gorm.io/gorm"
)

// RegistrySettingsRepository persists the single-row internal-registry config.
type RegistrySettingsRepository struct {
	db *gorm.DB
}

func NewRegistrySettingsRepository(db *gorm.DB) *RegistrySettingsRepository {
	return &RegistrySettingsRepository{db: db}
}

// Get returns the single settings row, or gorm.ErrRecordNotFound when unset.
func (r *RegistrySettingsRepository) Get() (*models.RegistrySettings, error) {
	var st models.RegistrySettings
	if err := r.db.Order("id ASC").First(&st).Error; err != nil {
		return nil, err
	}
	return &st, nil
}

// Upsert saves the single settings row (creating it on first save).
func (r *RegistrySettingsRepository) Upsert(st *models.RegistrySettings) error {
	return r.db.Save(st).Error
}
