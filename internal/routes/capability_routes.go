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
	"github.com/miabi-io/miabi/internal/dto"
	"github.com/miabi-io/miabi/internal/handlers"
)

// capabilityRoutes serve the grantable catalogue and the granted-app inventory.
func (r *Router) capabilityRoutes() []okapi.RouteDefinition {
	tag := okapi.GroupTag{Name: "Capabilities", Description: "Grantable Linux capabilities and host devices."}
	v1 := r.v1.Group("").WithTagInfo(tag)
	sys := r.v1.Group("/system").WithTagInfo(tag)
	authed := []okapi.Middleware{r.authenticate}
	admin := []okapi.Middleware{r.authenticate, r.systemAdmin}

	return []okapi.RouteDefinition{
		{
			Method:      http.MethodGet,
			Path:        "/capabilities",
			Group:       v1,
			Middlewares: authed,
			Handler:     r.h.capability.Catalog,
			Summary:     "Grantable Linux capabilities and host devices",
			Response:    &dto.Response[handlers.CapabilityCatalog]{},
		},
		{
			Method:      http.MethodGet,
			Path:        "/granted-applications",
			Group:       sys,
			Middlewares: admin,
			Handler:     r.h.capability.Inventory,
			Summary:     "Every application holding a capability or device grant (admin)",
			Response:    &dto.Response[[]handlers.GrantedApp]{},
		},
	}
}
