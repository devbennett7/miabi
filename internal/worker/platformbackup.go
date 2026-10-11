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

package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/jkaninda/logger"
	"github.com/miabi-io/miabi/internal/services/platformbackup"
)

// PlatformBackupHandler runs platform (control-plane) backup tasks.
type PlatformBackupHandler struct {
	svc *platformbackup.Service
}

func NewPlatformBackupHandler(svc *platformbackup.Service) *PlatformBackupHandler {
	return &PlatformBackupHandler{svc: svc}
}

// ProcessTask implements asynq.Handler for the platform-backup task.
func (h *PlatformBackupHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var p PlatformBackupPayload
	if err := json.Unmarshal(task.Payload(), &p); err != nil {
		return fmt.Errorf("bad platform backup payload: %w", err)
	}
	if err := h.svc.RunBackup(ctx, p.PlatformBackupID); err != nil {
		logger.Error("platform backup failed", "backup", p.PlatformBackupID, "error", err)
		return err
	}
	return nil
}
