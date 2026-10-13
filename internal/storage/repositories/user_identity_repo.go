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

type UserIdentityRepository struct {
	db *gorm.DB
}

func NewUserIdentityRepository(db *gorm.DB) *UserIdentityRepository {
	return &UserIdentityRepository{db: db}
}

func (r *UserIdentityRepository) Create(i *models.UserIdentity) error {
	return r.db.Create(i).Error
}

// FindBySubject returns the identity a provider's subject is bound to.
func (r *UserIdentityRepository) FindBySubject(providerID uint, subject string) (*models.UserIdentity, error) {
	var i models.UserIdentity
	if err := r.db.Where("provider_id = ? AND subject = ?", providerID, subject).First(&i).Error; err != nil {
		return nil, err
	}
	return &i, nil
}

// FindByUser returns the identity a user holds at a provider.
func (r *UserIdentityRepository) FindByUser(providerID, userID uint) (*models.UserIdentity, error) {
	var i models.UserIdentity
	if err := r.db.Where("provider_id = ? AND user_id = ?", providerID, userID).First(&i).Error; err != nil {
		return nil, err
	}
	return &i, nil
}

func (r *UserIdentityRepository) Delete(id uint) error {
	return r.db.Delete(&models.UserIdentity{}, id).Error
}
