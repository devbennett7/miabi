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

package components

// Component versions, bumped in the PR that changes the component: major when users must act (a
// migration, a manifest change, a removed option), minor for a feature, patch for a fix.
// scripts/components/check-component-versions.sh warns when a component's code changes without a bump here.
// The Web UI ships with Miabi and carries Miabi's own version.
const (
	GitOpsVersion         = "1.1.0"
	PipelineVersion       = "1.1.0"
	MarketplaceVersion    = "1.0.0"
	RegistryVersion       = "1.0.1"
	StorageClassesVersion = "1.0.0"
	AnalyticsVersion      = "1.0.0"
	HealthProbeVersion    = "1.0.0"
	SealedSecretsVersion  = "1.0.0"
)
