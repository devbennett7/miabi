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

package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UIDModel is an embeddable mixin giving a resource a stable UUIDv7 alongside its uint primary
// key: the PK stays the internal handle, the uid is the portable external one (Terraform id,
// GitOps/backup reference, non-enumerable public handle). Embed it anonymously.
type UIDModel struct {
	UID string `json:"uid" gorm:"type:uuid;uniqueIndex;not null;default:gen_random_uuid()"`
}

// BeforeCreate assigns a UUIDv7 when one was not supplied. A caller may preset
// UID (e.g. import/restore keeping an incoming identity); only an empty value is
// generated.
func (m *UIDModel) BeforeCreate(*gorm.DB) error {
	if m.UID == "" {
		v, err := uuid.NewV7()
		if err != nil {
			return err
		}
		m.UID = v.String()
	}
	return nil
}
