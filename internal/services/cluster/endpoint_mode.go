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
	"errors"

	"github.com/miabi-io/miabi/internal/models"
)

// ErrInvalidEndpointMode is returned for a service endpoint mode other than "vip" or "dnsrr".
var ErrInvalidEndpointMode = errors.New("service load balancing must be vip or dnsrr")

// ServiceEndpointMode is how a cluster's service apps are reached by name; id 0 is the default cluster.
func (s *Service) ServiceEndpointMode(clusterID uint) models.ServiceEndpointMode {
	if s.store == nil {
		return models.ServiceEndpointVIP
	}
	c, err := s.store.FindByID(clusterID)
	if err != nil {
		return models.ServiceEndpointVIP
	}
	return c.EndpointMode()
}

// SetEndpointModeListener is told when a cluster's service endpoint mode changed, so its running services follow.
func (s *Service) SetEndpointModeListener(fn func(ctx context.Context, clusterID uint, mode models.ServiceEndpointMode)) {
	s.endpointModeListener = fn
}

func (s *Service) endpointModeChanged(clusterID uint, mode models.ServiceEndpointMode) {
	if s.endpointModeListener != nil {
		go s.endpointModeListener(context.Background(), clusterID, mode)
	}
}
