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

package routes

import (
	"net/http"

	"github.com/jkaninda/okapi"
	"github.com/miabi-io/miabi/internal/middlewares"
	"github.com/miabi-io/miabi/internal/models"
)

func (r *Router) searchRoutes() []okapi.RouteDefinition {
	g := r.v1.Group("/workspaces").WithTagInfo(okapi.GroupTag{
		Name:        "Search",
		Description: "Cross-resource search within a workspace.",
	})
	return []okapi.RouteDefinition{
		{
			Method:      http.MethodGet,
			Path:        "/{workspace}/search",
			Group:       g,
			Middlewares: []okapi.Middleware{r.authenticate, r.scope, middlewares.RequireRole(models.WorkspaceRoleViewer)},
			Handler:     r.h.search.Search,
			Summary:     "Search applications, databases, routes and other resources in a workspace",
		},
	}
}
