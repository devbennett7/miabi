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

// AuditLog is an append-only record of a mutating action.
type AuditLog struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	ActorID     *uint          `json:"actor_id" gorm:"index"`
	WorkspaceID *uint          `json:"workspace_id" gorm:"index"`
	Action      string         `json:"action" gorm:"index;not null"`
	TargetType  string         `json:"target_type"`
	TargetID    string         `json:"target_id"`
	IPAddress   string         `json:"ip_address"`
	Metadata    map[string]any `json:"metadata,omitempty" gorm:"serializer:json"`
	CreatedAt   time.Time      `json:"created_at"`
}
