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

type RegistryRepository struct {
	db *gorm.DB
}

func NewRegistryRepository(db *gorm.DB) *RegistryRepository { return &RegistryRepository{db: db} }

func (r *RegistryRepository) Create(reg *models.Registry) error { return r.db.Create(reg).Error }
func (r *RegistryRepository) Update(reg *models.Registry) error { return r.db.Save(reg).Error }
func (r *RegistryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Registry{}, id).Error
}

func (r *RegistryRepository) FindInWorkspace(workspaceID, id uint) (*models.Registry, error) {
	var reg models.Registry
	if err := r.db.Where("id = ? AND workspace_id = ?", id, workspaceID).First(&reg).Error; err != nil {
		return nil, err
	}
	return &reg, nil
}

func (r *RegistryRepository) ListByWorkspace(workspaceID uint) ([]models.Registry, error) {
	var regs []models.Registry
	err := r.db.Where("workspace_id = ?", workspaceID).Order("created_at DESC").Find(&regs).Error
	return regs, err
}

func (r *RegistryRepository) ExistsByName(workspaceID uint, name string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Registry{}).
		Where("workspace_id = ? AND name = ?", workspaceID, name).Count(&count).Error
	return count > 0, err
}
