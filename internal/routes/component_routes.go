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
	"github.com/miabi-io/miabi/internal/components"
	"github.com/miabi-io/miabi/internal/dto"
)

// componentRoutes serve the independently versioned components shown on the About page.
func (r *Router) componentRoutes() []okapi.RouteDefinition {
	sys := r.v1.Group("/system").WithTagInfo(okapi.GroupTag{Name: "Components", Description: "Independently versioned parts of Miabi."})
	return []okapi.RouteDefinition{
		{
			Method:      http.MethodGet,
			Path:        "/components",
			Group:       sys,
			Middlewares: []okapi.Middleware{r.authenticate},
			Handler:     r.h.components.List,
			Summary:     "Component versions, stability and status (any signed-in user)",
			Response:    &dto.Response[[]components.Component]{},
		},
	}
}
