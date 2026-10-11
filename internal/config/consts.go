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

package config

const (
	// OfficialMarketplaceURL is Miabi's hosted template catalog and the default for
	// MIABI_MARKETPLACE_URL. Exported because it is also the fallback an unlicensed
	// install is held to when it configures a custom marketplace.
	OfficialMarketplaceURL = "https://marketplace.miabi.io"
	DefaultGomaImage       = "jkaninda/goma-gateway:" + DefaultGomaVersion

	// DefaultGomaVersion is the gateway version this build provisions and is
	// tested against. Keep it in step with GOMA_VERSION in the release workflow.
	DefaultGomaVersion = "1.0.0"

	marketplaceURL = OfficialMarketplaceURL
	miabiLogDir    = "/var/lib/miabi/logs"
)
