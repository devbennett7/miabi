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

package drift

import (
	"context"
	"errors"

	"github.com/miabi-io/miabi/internal/docker"
)

// Replaced reports whether a volume that exists under its own name is no longer the one Miabi recorded:
// its engine creation timestamp moved, so it was deleted and created again, and it holds none of the data
// the record describes. A volume Miabi never recorded a timestamp for is not replaced — most predate the
// column, and calling them all replaced would cry wolf about every one of them.
func Replaced(recorded, engineCreatedAt string) bool {
	return recorded != "" && engineCreatedAt != "" && engineCreatedAt != recorded
}

// VolumeState reports whether a volume still holds the data its record describes: "" when it does,
// ClassMissing when it is gone, ClassReplaced when it was deleted and recreated. An engine that cannot
// answer returns an error and no class — not knowing is never evidence of loss.
func VolumeState(ctx context.Context, dc docker.Client, name, recordedCreatedAt string) (string, error) {
	v, err := dc.InspectVolume(ctx, name)
	switch {
	case errors.Is(err, docker.ErrNotFound):
		return ClassMissing, nil
	case err != nil:
		return "", err
	case Replaced(recordedCreatedAt, v.CreatedAt):
		return ClassReplaced, nil
	default:
		return "", nil
	}
}
