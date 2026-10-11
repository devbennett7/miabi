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

import "github.com/hibiken/asynq"

// NewServer builds an asynq server with Miabi's queue priorities. embedded must be true only for the
// control-plane server's embedded worker: it is the one holding the agent tunnels (QueueNode) and the
// platform's full service graph (QueueControl). A standalone worker passes false.
func NewServer(redisAddr, redisPassword string, redisDB, concurrency int, embedded bool) *asynq.Server {
	queues := map[string]int{
		QueueDeploy:  6,
		QueueDefault: 3,
		QueueLow:     1,
	}
	if embedded {
		queues[QueueNode] = 6
		queues[QueueControl] = 3
	}
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr, Password: redisPassword, DB: redisDB},
		asynq.Config{
			Concurrency: concurrency,
			Queues:      queues,
		},
	)
}

// NewMux registers task handlers and returns the asynq mux.
func NewMux(deploy *DeployHandler, provision *ProvisionDBHandler, upgrade *UpgradeDBHandler, fanout *FanoutHandler, webhook *WebhookDeliverHandler, channel *ChannelSendHandler, job *JobHandler, volumeBackup *VolumeBackupHandler, pipeline *PipelineHandler, platformBackup *PlatformBackupHandler, wsBundle *WorkspaceBundleHandler, migration *LocationMigrationHandler) *asynq.ServeMux {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeDeploy, deploy.ProcessTask)
	mux.HandleFunc(TypeCanaryStep, deploy.ProcessCanaryStep)
	mux.HandleFunc(TypeProvisionDB, provision.ProcessTask)
	mux.HandleFunc(TypeUpgradeDB, upgrade.ProcessTask)
	mux.HandleFunc(TypeNotifyFanout, fanout.ProcessTask)
	mux.HandleFunc(TypeWebhookDeliver, webhook.ProcessTask)
	mux.HandleFunc(TypeNotifyChannel, channel.ProcessTask)
	mux.HandleFunc(TypeRunJob, job.ProcessTask)
	mux.HandleFunc(TypeVolumeBackup, volumeBackup.ProcessTask)
	mux.HandleFunc(TypeRunPipeline, pipeline.ProcessTask)
	mux.HandleFunc(TypePlatformBackup, platformBackup.ProcessTask)
	// Portable workspace bundles run only where the full service graph exists, so
	// a standalone worker passes nil here and leaves TypeWSBundle unregistered —
	// its server does not consume QueueControl either, so nothing is stranded.
	if wsBundle != nil {
		mux.HandleFunc(TypeWSBundle, wsBundle.ProcessTask)
	}
	if migration != nil {
		mux.HandleFunc(TypeMigration, migration.ProcessTask)
	}
	return mux
}
