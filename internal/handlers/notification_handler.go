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
	"github.com/miabi-io/miabi/internal/models"
	"github.com/miabi-io/miabi/internal/services/audit"
	"github.com/miabi-io/miabi/internal/services/notify"
)

type NotificationHandler struct {
	svc   *notify.Service
	audit *audit.Logger
}

func NewNotificationHandler(svc *notify.Service, auditLog *audit.Logger) *NotificationHandler {
	return &NotificationHandler{svc: svc, audit: auditLog}
}

type CreateChannelRequest struct {
	Body struct {
		Type    string   `json:"type" required:"true" enum:"telegram,slack,discord"`
		Name    string   `json:"name" required:"true"`
		Events  []string `json:"events" required:"true" minItems:"1"`
		Enabled bool     `json:"enabled"`
		// Telegram.
		BotToken string `json:"bot_token"`
		ChatID   string `json:"chat_id"`
		// Slack / Discord.
		WebhookURL string `json:"webhook_url"`
	} `json:"body"`
}

type UpdateChannelRequest struct {
	Body struct {
		Name       string   `json:"name"`
		Events     []string `json:"events"`
		Enabled    bool     `json:"enabled"`
		BotToken   string   `json:"bot_token"`
		ChatID     string   `json:"chat_id"`
		WebhookURL string   `json:"webhook_url"`
	} `json:"body"`
}

func (h *NotificationHandler) List(c *okapi.Context) error {
	chs, err := h.svc.List(middlewares.WorkspaceID(c))
	if err != nil {
		return c.AbortInternalServerError("failed to list notification channels", err)
	}
	return ok(c, chs)
}

func (h *NotificationHandler) Create(c *okapi.Context, req *CreateChannelRequest) error {
	wsID := middlewares.WorkspaceID(c)
	ch, err := h.svc.Create(wsID, notify.Input{
		Type:       models.NotificationChannelType(req.Body.Type),
		Name:       req.Body.Name,
		BotToken:   req.Body.BotToken,
		ChatID:     req.Body.ChatID,
		WebhookURL: req.Body.WebhookURL,
		Events:     req.Body.Events,
		Enabled:    req.Body.Enabled,
	})
	if err != nil {
		return h.mapErr(c, err)
	}
	h.record(c, wsID, "notification_channel.create", ch.ID)
	return created(c, ch)
}

func (h *NotificationHandler) Get(c *okapi.Context) error {
	id, err := h.id(c)
	if err != nil {
		return c.AbortBadRequest("invalid channel id")
	}
	ch, err := h.svc.Get(middlewares.WorkspaceID(c), id)
	if err != nil {
		return c.AbortNotFound("notification channel not found")
	}
	return ok(c, ch)
}

func (h *NotificationHandler) Update(c *okapi.Context, req *UpdateChannelRequest) error {
	id, err := h.id(c)
	if err != nil {
		return c.AbortBadRequest("invalid channel id")
	}
	wsID := middlewares.WorkspaceID(c)
	ch, err := h.svc.Update(wsID, id, notify.Input{
		Name:       req.Body.Name,
		BotToken:   req.Body.BotToken,
		ChatID:     req.Body.ChatID,
		WebhookURL: req.Body.WebhookURL,
		Events:     req.Body.Events,
		Enabled:    req.Body.Enabled,
	})
	if err != nil {
		return h.mapErr(c, err)
	}
	h.record(c, wsID, "notification_channel.update", ch.ID)
	return ok(c, ch)
}

func (h *NotificationHandler) Delete(c *okapi.Context) error {
	id, err := h.id(c)
	if err != nil {
		return c.AbortBadRequest("invalid channel id")
	}
	wsID := middlewares.WorkspaceID(c)
	if err := h.svc.Delete(wsID, id); err != nil {
		return h.mapErr(c, err)
	}
	h.record(c, wsID, "notification_channel.delete", id)
	return message(c, "notification channel deleted")
}

func (h *NotificationHandler) Test(c *okapi.Context) error {
	id, err := h.id(c)
	if err != nil {
		return c.AbortBadRequest("invalid channel id")
	}
	if err := h.svc.Test(c.Request().Context(), middlewares.WorkspaceID(c), id); err != nil {
		if errors.Is(err, notify.ErrNotFound) {
			return c.AbortNotFound("notification channel not found")
		}
		return c.AbortWithError(400, err)
	}
	return message(c, "test notification sent")
}

func (h *NotificationHandler) id(c *okapi.Context) (uint, error) {
	id, err := strconv.Atoi(c.Param("channelID"))
	if err != nil || id <= 0 {
		return 0, errors.New("invalid channel id")
	}
	return uint(id), nil
}

func (h *NotificationHandler) record(c *okapi.Context, wsID uint, action string, id uint) {
	actor := middlewares.UserID(c)
	h.audit.Record(audit.Entry{ActorID: &actor, WorkspaceID: &wsID, Action: action, TargetType: "notification_channel", TargetID: strconv.Itoa(int(id)), IP: c.RealIP()})
}

func (h *NotificationHandler) mapErr(c *okapi.Context, err error) error {
	switch {
	case errors.Is(err, notify.ErrNameRequired), errors.Is(err, notify.ErrUnsupportedType),
		errors.Is(err, notify.ErrTokenRequired), errors.Is(err, notify.ErrChatIDRequired),
		errors.Is(err, notify.ErrWebhookURLRequired), errors.Is(err, notify.ErrInvalidEvent):
		return c.AbortBadRequest(err.Error())
	case errors.Is(err, notify.ErrNotFound):
		return c.AbortNotFound("notification channel not found")
	default:
		return c.AbortInternalServerError("notification channel operation failed", err)
	}
}
