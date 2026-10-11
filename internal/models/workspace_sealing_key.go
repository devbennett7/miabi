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

import "time"

// WorkspaceSealingKey is a workspace's key pair for sealed secrets: values encrypted to PublicKey can be
// committed to git, and only Miabi, holding the private half, opens them at apply time. Rotation adds a new
// Active version; older versions are kept so values already in git keep opening.
type WorkspaceSealingKey struct {
	ID          uint `json:"-" gorm:"primaryKey"`
	WorkspaceID uint `json:"-" gorm:"index:idx_wsseal_ws_ver,unique;not null"`
	Version     int  `json:"version" gorm:"index:idx_wsseal_ws_ver,unique;not null"`
	// PublicKey is the age recipient ("age1…") that clients seal to.
	PublicKey string `json:"public_key" gorm:"type:text;not null"`
	// IdentityEnc is the private age identity wrapped under the master KEK, not the workspace DEK, so a DEK
	// rotation never touches it and a DR restore of the database brings it back with the master key.
	IdentityEnc string    `json:"-" gorm:"type:text;not null"`
	Active      bool      `json:"active" gorm:"not null;default:true;index"`
	CreatedAt   time.Time `json:"created_at"`
}
