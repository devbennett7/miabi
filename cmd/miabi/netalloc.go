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

package main

import (
	"github.com/jkaninda/logger"
	"github.com/miabi-io/miabi/internal/config"
	"github.com/miabi-io/miabi/internal/services/netalloc"
	"github.com/miabi-io/miabi/internal/storage/repositories"
	"gorm.io/gorm"
)

// newSubnetAllocator builds the subnet allocator for the worker's deploy/job handlers.
// Returns nil on a config error, so they fall back to Docker's default address pool.
func newSubnetAllocator(cfg *config.Config, db *gorm.DB) *netalloc.Service {
	a, err := netalloc.NewService(repositories.NewNetworkAllocationRepository(db), cfg.NetworkPoolCIDR, cfg.NetworkSubnetPrefix)
	if err != nil {
		logger.Warn("subnet allocator disabled in worker; falling back to Docker's default address pool", "error", err)
		return nil
	}
	return a
}
