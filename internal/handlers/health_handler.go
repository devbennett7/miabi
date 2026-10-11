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
	"context"
	"time"

	"github.com/jkaninda/okapi"
	"github.com/miabi-io/miabi/internal/docker"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const statusNotReady = "not ready"

// HealthHandler serves liveness and readiness probes.
type HealthHandler struct {
	db     *gorm.DB
	redis  *redis.Client
	docker docker.Client
}

func NewHealthHandler(db *gorm.DB, redis *redis.Client, dockerClient docker.Client) *HealthHandler {
	return &HealthHandler{db: db, redis: redis, docker: dockerClient}
}

type HealthResponse struct {
	Status string `json:"status" example:"ok"`
}

type ReadyResponse struct {
	Status   string `json:"status" example:"ready"`
	Database string `json:"database" example:"ok"`
	Redis    string `json:"redis" example:"ok"`
	Docker   string `json:"docker" example:"ok"`
}

// Healthz is a lightweight liveness probe.
func (h *HealthHandler) Healthz(c *okapi.Context) error {
	return c.OK(HealthResponse{Status: "ok"})
}

// Readyz checks that all dependencies are reachable.
func (h *HealthHandler) Readyz(c *okapi.Context) error {
	resp := ReadyResponse{Status: "ready", Database: "ok", Redis: "ok", Docker: "ok"}

	if sqlDB, err := h.db.DB(); err != nil {
		resp.Status = statusNotReady
		resp.Database = err.Error()
	} else {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
		defer cancel()
		if err := sqlDB.PingContext(ctx); err != nil {
			resp.Status = statusNotReady
			resp.Database = err.Error()
		}
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
	defer cancel()
	if err := h.redis.Ping(ctx).Err(); err != nil {
		resp.Status = statusNotReady
		resp.Redis = err.Error()
	}

	dctx, dcancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
	defer dcancel()
	if err := h.docker.Ping(dctx); err != nil {
		resp.Status = statusNotReady
		resp.Docker = err.Error()
	}

	if resp.Status != "ready" {
		return c.JSON(503, resp)
	}
	return c.OK(resp)
}
