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

package worker

// ImageRefResolver resolves a platform-image catalog key to its effective ref (admin override ->
// config default -> built-in). Kept as an interface so the worker doesn't import the settings
// stack. Builds themselves run on the runner; the control plane only passes it the builder image.
type ImageRefResolver interface {
	Ref(key string) string
}

// BuildResult is what a runner build reports for provenance: the pushed image
// digest, its size when known, and which runner built it.
type BuildResult struct {
	Digest string
	Size   int64
	Runner string
}
