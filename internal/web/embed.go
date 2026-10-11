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

// Package web embeds the built frontend (the Vue SPA) into the Go binary so a single `miabi`
// executable serves both the API and the UI. `dist/` is a build artifact: only a .gitkeep
// placeholder is committed, so `go build` works even when the UI has not been built.
package web

import "embed"

// Assets holds the built SPA under a top-level "dist/" directory, served via Okapi's WebFS
// with WebConfig{Root: "dist"}. The `all:` prefix keeps dotted files (the .gitkeep
// placeholder, Vite's occasional dotfiles) in the embed.
//
//go:embed all:dist
var Assets embed.FS
