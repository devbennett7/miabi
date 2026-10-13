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

package application

import (
	"context"
	"testing"

	"github.com/miabi-io/miabi/internal/docker"
	"github.com/miabi-io/miabi/internal/models"
)

type listingEngine struct {
	docker.Client
	services []docker.ServiceStatus
	calls    int
}

func (e *listingEngine) ServiceList(context.Context) ([]docker.ServiceStatus, error) {
	e.calls++
	return e.services, nil
}

type swarmWith struct{ mgr *listingEngine }

func (c swarmWith) IsSwarm(uint) bool                                    { return true }
func (c swarmWith) Manager(context.Context, uint) (docker.Client, error) { return c.mgr, nil }
func (c swarmWith) ClusterOfServer(uint) uint                            { return models.DefaultClusterID }

// A list asks the swarm once, not once per app, and only fills the apps it reports.
func TestAnnotateReplicasListsServicesOncePerCluster(t *testing.T) {
	mgr := &listingEngine{services: []docker.ServiceStatus{{Name: "mb-app-a-1", Replicas: 3, RunningTasks: 2}}}
	s := &Service{cluster: swarmWith{mgr: mgr}}
	apps := []models.Application{
		{ID: 1, Alias: "mb-app-a-1", RuntimeKind: models.RuntimeService},
		{ID: 2, Alias: "mb-app-b-2", RuntimeKind: models.RuntimeService},
		{ID: 3, Alias: "mb-app-c-3", RuntimeKind: models.RuntimeContainer},
	}

	s.annotateReplicas(apps)

	if mgr.calls != 1 {
		t.Fatalf("ServiceList called %d times, want once", mgr.calls)
	}
	if r, d := apps[0].RunningReplicas, apps[0].DesiredReplicas; r == nil || d == nil || *r != 2 || *d != 3 {
		t.Fatalf("listed service counts = %v/%v, want 2/3", r, d)
	}
	if apps[1].RunningReplicas != nil || apps[2].RunningReplicas != nil {
		t.Fatal("counts were set for an unlisted service or a container app")
	}
}
