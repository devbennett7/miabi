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

// DNSProviderRepository persists workspace-scoped DNS provider connections.
type DNSProviderRepository struct {
	db *gorm.DB
}

func NewDNSProviderRepository(db *gorm.DB) *DNSProviderRepository {
	return &DNSProviderRepository{db: db}
}

func (r *DNSProviderRepository) Create(p *models.DNSProvider) error { return r.db.Create(p).Error }
func (r *DNSProviderRepository) Update(p *models.DNSProvider) error { return r.db.Save(p).Error }

func (r *DNSProviderRepository) FindInWorkspace(workspaceID, id uint) (*models.DNSProvider, error) {
	var p models.DNSProvider
	if err := r.db.Where("id = ? AND workspace_id = ?", id, workspaceID).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// FindByID looks up a provider by id alone (used by the reconcile sweep to update
// status; it works from the record ledger and has no workspace in hand).
func (r *DNSProviderRepository) FindByID(id uint) (*models.DNSProvider, error) {
	var p models.DNSProvider
	if err := r.db.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *DNSProviderRepository) ListByWorkspace(workspaceID uint) ([]models.DNSProvider, error) {
	var ps []models.DNSProvider
	err := r.db.Where("workspace_id = ?", workspaceID).Order("name ASC").Find(&ps).Error
	return ps, err
}

func (r *DNSProviderRepository) CountByWorkspace(workspaceID uint) (int64, error) {
	var n int64
	err := r.db.Model(&models.DNSProvider{}).Where("workspace_id = ?", workspaceID).Count(&n).Error
	return n, err
}

func (r *DNSProviderRepository) ExistsByName(workspaceID uint, name string) (bool, error) {
	var n int64
	err := r.db.Model(&models.DNSProvider{}).
		Where("workspace_id = ? AND name = ?", workspaceID, name).Count(&n).Error
	return n > 0, err
}

func (r *DNSProviderRepository) Delete(id uint) error {
	return r.db.Delete(&models.DNSProvider{}, id).Error
}
