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

package middlewares

import (
	"strconv"
	"strings"

	"github.com/jkaninda/okapi"
)

// AppIDResolver turns an application's portable uid into its numeric id. Satisfied by the
// application repository.
type AppIDResolver interface {
	IDByUID(uid string) (uint, error)
}

// AppRefResolver resolves an app reference (id, uid or name) within a workspace. The application
// repository satisfies it; a resolver without it accepts ids and uids only.
type AppRefResolver interface {
	IDByRef(workspaceID uint, ref string) (uint, error)
}

func confineEphemeralKey(c *okapi.Context, boundApp uint, apps AppIDResolver, workspaceID uint) error {
	if boundApp == 0 {
		return c.AbortForbidden("this job credential cannot be used on the API")
	}
	ref := strings.TrimSpace(c.Param("appID"))
	if ref == "" {
		return errJobKeyScope(c)
	}
	got, ok := resolveAppRef(ref, apps)
	// A job key is workspace-bound, so a name resolves against its own workspace.
	if r, isRef := apps.(AppRefResolver); isRef && workspaceID != 0 {
		if id, err := r.IDByRef(workspaceID, ref); err == nil {
			got, ok = id, true
		}
	}
	if !ok || got != boundApp {
		return errJobKeyScope(c)
	}
	return nil
}

func errJobKeyScope(c *okapi.Context) error {
	return c.AbortForbidden("this job credential may only act on the application it was issued for")
}

func resolveAppRef(ref string, apps AppIDResolver) (uint, bool) {
	if id, err := strconv.Atoi(ref); err == nil {
		if id <= 0 {
			return 0, false
		}
		return uint(id), true
	}
	if apps == nil {
		return 0, false
	}
	id, err := apps.IDByUID(ref)
	if err != nil || id == 0 {
		return 0, false
	}
	return id, true
}
