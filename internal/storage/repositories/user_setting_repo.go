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
	"errors"

	"github.com/miabi-io/miabi/internal/models"
	"gorm.io/gorm"
)

type UserSettingRepository struct {
	db *gorm.DB
}

func NewUserSettingRepository(db *gorm.DB) *UserSettingRepository {
	return &UserSettingRepository{db: db}
}

// Get returns a user's preferences, or the defaults when they have never saved any.
// A missing row is not an error: most users never open Preferences, and the caller
// wants a usable value rather than a branch.
func (r *UserSettingRepository) Get(userID uint) (*models.UserSetting, error) {
	var s models.UserSetting
	err := r.db.Where("user_id = ?", userID).First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		def := models.DefaultUserSetting()
		def.UserID = userID
		return &def, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// Save upserts a user's preferences. The row is created on first write, so a user
// only occupies a row once they have expressed a preference.
func (r *UserSettingRepository) Save(s *models.UserSetting) error {
	var existing models.UserSetting
	err := r.db.Where("user_id = ?", s.UserID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.ID = 0
		return r.db.Create(s).Error
	}
	if err != nil {
		return err
	}
	s.ID = existing.ID
	return r.db.Save(s).Error
}

// DeleteByUser removes a user's preferences, for account deletion.
func (r *UserSettingRepository) DeleteByUser(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.UserSetting{}).Error
}
