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
	"github.com/miabi-io/miabi/internal/services/volumebackup"
)

// VolumeBackupHandler runs volume backup tasks.
type VolumeBackupHandler struct {
	svc *volumebackup.Service
}

func NewVolumeBackupHandler(svc *volumebackup.Service) *VolumeBackupHandler {
	return &VolumeBackupHandler{svc: svc}
}

// ProcessTask implements asynq.Handler for the volume-backup task.
func (h *VolumeBackupHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var p VolumeBackupPayload
	if err := json.Unmarshal(task.Payload(), &p); err != nil {
		return fmt.Errorf("bad volume backup payload: %w", err)
	}
	if err := h.svc.RunBackup(ctx, p.VolumeBackupID); err != nil {
		logger.Error("volume backup failed", "backup", p.VolumeBackupID, "error", err)
		return err
	}
	return nil
}
