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

package handlers

import (
	"errors"
	"strconv"

	"github.com/jkaninda/okapi"
	"github.com/miabi-io/miabi/internal/middlewares"
	"github.com/miabi-io/miabi/internal/services/audit"
	"github.com/miabi-io/miabi/internal/services/release"
)

// ReleaseHandler exposes the workspace release catalog, approvals, and promotion.
type ReleaseHandler struct {
	svc   *release.Service
	audit *audit.Logger
}

func NewReleaseHandler(svc *release.Service, auditLog *audit.Logger) *ReleaseHandler {
	return &ReleaseHandler{svc: svc, audit: auditLog}
}

type ApproveReleaseRequest struct {
	Body struct {
		EnvironmentID *uint  `json:"environment_id"`
		Approved      bool   `json:"approved" default:"true"`
		Comment       string `json:"comment"`
	} `json:"body"`
}

type PromoteReleaseRequest struct {
	Body struct {
		EnvironmentID uint `json:"environment_id" required:"true"`
	} `json:"body"`
}

func (h *ReleaseHandler) List(c *okapi.Context) error {
	page, size, offset := normalizePageParams(queryInt(c, "page", 0), queryInt(c, "size", 20))
	out, total, err := h.svc.ListPaged(middlewares.WorkspaceID(c), size, offset)
	if err != nil {
		return c.AbortInternalServerError("failed to list releases", err)
	}
	return paginated(c, out, total, page, size)
}

func (h *ReleaseHandler) Approvals(c *okapi.Context) error {
	id, err := uintParam(c, "releaseID")
	if err != nil {
		return c.AbortBadRequest("invalid release id")
	}
	status, err := h.svc.Approvals(middlewares.WorkspaceID(c), id)
	if err != nil {
		return h.mapErr(c, err)
	}
	return ok(c, status)
}

func (h *ReleaseHandler) Approve(c *okapi.Context, req *ApproveReleaseRequest) error {
	id, err := uintParam(c, "releaseID")
	if err != nil {
		return c.AbortBadRequest("invalid release id")
	}
	wsID := middlewares.WorkspaceID(c)
	actor := middlewares.UserID(c)
	if err := h.svc.Approve(wsID, id, req.Body.EnvironmentID, actor, req.Body.Approved, req.Body.Comment); err != nil {
		return h.mapErr(c, err)
	}
	h.record(c, wsID, "release.approve", id)
	return message(c, "approval recorded")
}

func (h *ReleaseHandler) Promote(c *okapi.Context, req *PromoteReleaseRequest) error {
	id, err := uintParam(c, "releaseID")
	if err != nil {
		return c.AbortBadRequest("invalid release id")
	}
	wsID := middlewares.WorkspaceID(c)
	actor := middlewares.UserID(c)
	dep, err := h.svc.Promote(wsID, id, req.Body.EnvironmentID, actor)
	if err != nil {
		return h.mapErr(c, err)
	}
	h.record(c, wsID, "release.promote", id)
	return created(c, dep)
}

func (h *ReleaseHandler) record(c *okapi.Context, wsID uint, action string, id uint) {
	actor := middlewares.UserID(c)
	h.audit.Record(audit.Entry{ActorID: &actor, WorkspaceID: &wsID, Action: action,
		TargetType: "release", TargetID: strconv.Itoa(int(id)), IP: c.RealIP()})
}

func (h *ReleaseHandler) mapErr(c *okapi.Context, err error) error {
	switch {
	case errors.Is(err, release.ErrNotFound):
		return c.AbortNotFound("release not found")
	case errors.Is(err, release.ErrEnvNotFound):
		return c.AbortNotFound("environment not found")
	case errors.Is(err, release.ErrApprovalsNeeded):
		return c.AbortForbidden(err.Error())
	default:
		return c.AbortInternalServerError("release operation failed", err)
	}
}
