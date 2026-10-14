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

// Components an UpdateStatus row may track.
const (
	UpdateComponentMiabi   = "miabi"
	UpdateComponentAgent   = "agent"
	UpdateComponentRunner  = "runner"
	UpdateComponentGateway = "gateway"
)

// UpdateStatus caches the result of a daily release check.
type UpdateStatus struct {
	ID            uint       `json:"-" gorm:"primaryKey"`
	Component     string     `json:"component" gorm:"uniqueIndex;not null;default:miabi"`
	LatestVersion string     `json:"latest_version"`
	ReleaseURL    string     `json:"release_url"`
	PublishedAt   *time.Time `json:"published_at"`
	// Reason is why the release was published: bugfix, feature, security or other.
	Reason           string     `json:"reason,omitempty"`
	ETag             string     `json:"-"`
	CheckedVersion   string     `json:"-"`
	CheckedAt        *time.Time `json:"checked_at"`
	LastError        string     `json:"last_error,omitempty"`
	DismissedVersion string     `json:"dismissed_version,omitempty"`
	UpdatedAt        time.Time  `json:"-"`
}
