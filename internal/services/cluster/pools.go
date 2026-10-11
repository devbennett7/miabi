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

package cluster

import (
	"context"

	"github.com/jkaninda/logger"
	"github.com/miabi-io/miabi/internal/models"
)

// SyncPoolLabel mirrors a node's pool onto its swarm node, where service placement constraints read it. A no-op
// for a node outside an active swarm.
func (s *Service) SyncPoolLabel(ctx context.Context, srv *models.Server) error {
	if srv == nil || srv.SwarmNodeID == "" || !s.IsSwarm(srv.ClusterID) {
		return nil
	}
	mgr, err := s.Manager(ctx, srv.ClusterID)
	if err != nil {
		return err
	}
	return mgr.SwarmNodeSetLabel(ctx, srv.SwarmNodeID, models.PoolLabel, models.PoolOf(srv))
}

// reassertPoolLabels labels every pooled swarm member, covering a node that joined after its pool was set or a
// pool set while its swarm was unreachable.
func (s *Service) reassertPoolLabels(ctx context.Context) {
	servers, err := s.nodes.List(ctx)
	if err != nil {
		return
	}
	for i := range servers {
		if models.PoolOf(&servers[i]) == "" {
			continue
		}
		if err := s.SyncPoolLabel(ctx, &servers[i]); err != nil {
			logger.Warn("failed to label a swarm node with its pool", "node", servers[i].Name, "error", err)
		}
	}
}
