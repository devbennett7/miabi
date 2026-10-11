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

type StackEnvVarRepository struct {
	db *gorm.DB
}

func NewStackEnvVarRepository(db *gorm.DB) *StackEnvVarRepository {
	return &StackEnvVarRepository{db: db}
}

func (r *StackEnvVarRepository) ListByStack(stackID uint) ([]models.StackEnvVar, error) {
	var vars []models.StackEnvVar
	err := r.db.Where("stack_id = ?", stackID).Order("key ASC").Find(&vars).Error
	return vars, err
}

func (r *StackEnvVarRepository) Upsert(v *models.StackEnvVar) error {
	var existing models.StackEnvVar
	err := r.db.Where("stack_id = ? AND key = ?", v.StackID, v.Key).First(&existing).Error
	if err == nil {
		existing.Value = v.Value
		existing.IsSecret = v.IsSecret
		return r.db.Save(&existing).Error
	}
	return r.db.Create(v).Error
}

func (r *StackEnvVarRepository) Delete(stackID uint, key string) error {
	return r.db.Where("stack_id = ? AND key = ?", stackID, key).Delete(&models.StackEnvVar{}).Error
}
