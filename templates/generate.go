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

package templates

// The embedded official catalog is a vendored snapshot of the marketplace
// registry's official/ set — the offline floor. It is GENERATED, not hand-kept:
// `go generate ./templates` rebuilds the slug trees, index.yaml, and the
// //go:embed list in embed.go from the marketplace's published export.json.
//
// Source defaults to the sibling repo's committed bundle; for a release, vendor
// from the published bundle instead, e.g.:
//
//	go run ./gen -source https://github.com/miabi-io/marketplace/releases/latest/download/export.json -out .
//
//go:generate go run ./gen -source ../../marketplace/export.json -out .
