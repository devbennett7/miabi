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

// ACMEAccountRepository persists the platform ACME account(s), keyed by CA URL.
type ACMEAccountRepository struct {
	db *gorm.DB
}

func NewACMEAccountRepository(db *gorm.DB) *ACMEAccountRepository {
	return &ACMEAccountRepository{db: db}
}

func (r *ACMEAccountRepository) Create(a *models.ACMEAccount) error { return r.db.Create(a).Error }

// FindByCA returns the account for a CA directory URL, or ErrRecordNotFound.
func (r *ACMEAccountRepository) FindByCA(caDirURL string) (*models.ACMEAccount, error) {
	var a models.ACMEAccount
	if err := r.db.Where("ca_dir_url = ?", caDirURL).First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}
