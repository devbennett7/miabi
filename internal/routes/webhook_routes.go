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
	"github.com/miabi-io/miabi/internal/handlers"
	"github.com/miabi-io/miabi/internal/middlewares"
	"github.com/miabi-io/miabi/internal/models"
)

// webhookRoutes registers workspace webhooks. Viewer reads; Admin creates, updates, deletes and tests.
func (r *Router) webhookRoutes() []okapi.RouteDefinition {
	g := r.v1.Group("/workspaces").WithTagInfo(okapi.GroupTag{Name: "Webhooks", Description: "Outbound HTTP webhooks for application events."})
	scoped := func(min models.WorkspaceRole) []okapi.Middleware {
		return []okapi.Middleware{r.authenticate, r.scope, middlewares.RequireRole(min)}
	}
	const base = "/{workspace}/webhooks"

	return []okapi.RouteDefinition{
		{
			Method:      http.MethodGet,
			Path:        base,
			Group:       g,
			Middlewares: scoped(models.WorkspaceRoleViewer),
			Handler:     r.h.webhook.List,
			Summary:     "List webhooks",
		},
		{
			Method:      http.MethodPost,
			Path:        base,
			Group:       g,
			Middlewares: scoped(models.WorkspaceRoleAdmin),
			Handler:     okapi.H(r.h.webhook.Create),
			Summary:     "Create a webhook",
			Request:     &handlers.CreateWebhookRequest{},
		},
		{
			Method:      http.MethodGet,
			Path:        base + "/{webhookID}",
			Group:       g,
			Middlewares: scoped(models.WorkspaceRoleViewer),
			Handler:     r.h.webhook.Get,
			Summary:     "Get a webhook",
		},
		{
			Method:      http.MethodPatch,
			Path:        base + "/{webhookID}",
			Group:       g,
			Middlewares: scoped(models.WorkspaceRoleAdmin),
			Handler:     okapi.H(r.h.webhook.Update),
			Summary:     "Update a webhook",
			Request:     &handlers.UpdateWebhookRequest{},
		},
		{
			Method:      http.MethodDelete,
			Path:        base + "/{webhookID}",
			Group:       g,
			Middlewares: scoped(models.WorkspaceRoleAdmin),
			Handler:     r.h.webhook.Delete,
			Summary:     "Delete a webhook",
		},
		{
			Method:      http.MethodPost,
			Path:        base + "/{webhookID}/test",
			Group:       g,
			Middlewares: scoped(models.WorkspaceRoleAdmin),
			Handler:     r.h.webhook.Test,
			Summary:     "Send a test delivery",
		},
		{
			Method:      http.MethodGet,
			Path:        base + "/{webhookID}/deliveries",
			Group:       g,
			Middlewares: scoped(models.WorkspaceRoleViewer),
			Handler:     r.h.webhook.Deliveries,
			Summary:     "List recent deliveries",
		},
		{
			Method:      http.MethodPost,
			Path:        base + "/{webhookID}/deliveries/{deliveryID}/redeliver",
			Group:       g,
			Middlewares: scoped(models.WorkspaceRoleAdmin),
			Handler:     r.h.webhook.Redeliver,
			Summary:     "Re-send a past delivery",
		},
	}
}
