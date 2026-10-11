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

package dr

import (
	"fmt"
	"path"
	"strings"
	"time"
)

// RefPrefix marks a platform recovery point. Restore matches on it so a
// mistyped argument fails fast instead of listing a whole bucket.
const RefPrefix = "mbdr_"

// IdentityExt is the sealed identity envelope's file extension.
const IdentityExt = ".mbid"

// NewRef builds a recovery point's stable name from the install it came from and
// the moment it started: "mbdr_<install-id>_20260731T020000Z".
func NewRef(installID string, at time.Time) string {
	id := strings.TrimSpace(installID)
	if id == "" {
		id = "unknown"
	}
	return fmt.Sprintf("%s%s_%s", RefPrefix, id, at.UTC().Format("20060102T150405Z"))
}

// IsRef reports whether s looks like a recovery point ref.
func IsRef(s string) bool { return strings.HasPrefix(strings.TrimSpace(s), RefPrefix) }

// IdentityObject is the identity envelope's object name within a bucket prefix.
// It sits beside the artifacts of its own set, named for it, so an operator
// listing a bucket can see which envelope belongs to which recovery point.
func IdentityObject(prefix, ref string) string {
	name := "identity-" + ref + IdentityExt
	if p := strings.Trim(strings.TrimSpace(prefix), "/"); p != "" {
		return path.Join(p, name)
	}
	return name
}

// RefFromIdentityObject recovers the ref from an identity object key, or "" when
// the key is not one. Lets `miabi restore` discover recovery points in a bucket
// without a database to list them from.
func RefFromIdentityObject(key string) string {
	base := path.Base(strings.TrimSpace(key))
	if !strings.HasPrefix(base, "identity-") || !strings.HasSuffix(base, IdentityExt) {
		return ""
	}
	ref := strings.TrimSuffix(strings.TrimPrefix(base, "identity-"), IdentityExt)
	if !IsRef(ref) {
		return ""
	}
	return ref
}
