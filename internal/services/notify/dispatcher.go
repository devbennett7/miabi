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

// Package notify turns persisted application events into outbound
// notifications: it fans an event out to a workspace's webhooks and channels,
// and renders/sends channel messages.
package notify

import (
	"github.com/jkaninda/logger"
	"github.com/miabi-io/miabi/internal/models"
)

// Enqueuer schedules the fan-out of a recorded event. Implemented by
// worker.Producer; defined here so this package does not import worker (which
// in turn imports this package for the sender registry).
type Enqueuer interface {
	EnqueueNotifyFanout(appEventID uint) error
}

// Dispatcher implements events.Notifier. It enqueues a single fan-out task per
// event; channel resolution and delivery happen in the worker.
type Dispatcher struct {
	enqueuer Enqueuer
}

func NewDispatcher(enqueuer Enqueuer) *Dispatcher { return &Dispatcher{enqueuer: enqueuer} }

// OnEvent schedules fan-out for a persisted, notifiable event. Best-effort: a
// failure is logged but never propagated to the recorder.
func (d *Dispatcher) OnEvent(e *models.AppEvent) {
	if d == nil || d.enqueuer == nil || e == nil || e.ID == 0 {
		return
	}
	if err := d.enqueuer.EnqueueNotifyFanout(e.ID); err != nil {
		logger.Error("failed to enqueue notification fan-out", "event", e.ID, "type", e.Type, "error", err)
	}
}
