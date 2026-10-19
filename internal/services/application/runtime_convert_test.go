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
	"errors"
	"strings"
	"testing"

	"github.com/miabi-io/miabi/internal/models"
)

func hasMsg(list []string, sub string) bool {
	for _, s := range list {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func deployed(app *models.Application) *models.Application {
	id := uint(1)
	app.CurrentReleaseID = &id
	return app
}

func TestPreflightContainerToService(t *testing.T) {
	s := &Service{cluster: clusterOn{}}
	app := deployed(&models.Application{RuntimeKind: models.RuntimeContainer, Replicas: 1})
	p := s.PreflightRuntime(app, RuntimeChange{Runtime: models.RuntimeService, Replicas: 3})
	if len(p.Blockers) != 0 {
		t.Fatalf("unexpected blockers: %v", p.Blockers)
	}
	if p.Replicas != 3 || p.Strategy != models.DeployRolling || !p.Redeploy {
		t.Errorf("preflight = %+v, want 3 replicas, rolling, redeploy", p)
	}
	if app.RuntimeKind != models.RuntimeContainer || app.Replicas != 1 {
		t.Error("a preflight must not change the app")
	}
}

func TestPreflightServiceToContainerDropsReplicas(t *testing.T) {
	s := &Service{cluster: clusterOn{}}
	app := deployed(&models.Application{RuntimeKind: models.RuntimeService, Replicas: 4})
	p := s.PreflightRuntime(app, RuntimeChange{Runtime: models.RuntimeContainer})
	if len(p.Blockers) != 0 || p.Replicas != 1 || !hasMsg(p.Warnings, "replicas drop from 4 to 1") {
		t.Fatalf("preflight = %+v", p)
	}
}

func TestPreflightBlockers(t *testing.T) {
	canary := uint(9)
	cases := []struct {
		name string
		s    *Service
		app  *models.Application
		to   models.RuntimeKind
		want string
	}{
		{"same runtime", &Service{cluster: clusterOn{}}, &models.Application{RuntimeKind: models.RuntimeContainer}, models.RuntimeContainer, "already runs"},
		{"cluster off", &Service{cluster: clusterOff{}}, &models.Application{RuntimeKind: models.RuntimeContainer}, models.RuntimeService, "cluster mode"},
		{"manifest", &Service{cluster: clusterOn{}}, &models.Application{RuntimeKind: models.RuntimeContainer,
			Metadata: map[string]string{models.MetaManagedBy: models.ManagedByGitOps}}, models.RuntimeService, "manifest"},
		{"canary", &Service{cluster: clusterOn{}}, &models.Application{RuntimeKind: models.RuntimeContainer, CanaryReleaseID: &canary}, models.RuntimeService, "canary"},
		{"gpu", &Service{cluster: clusterOn{}}, &models.Application{RuntimeKind: models.RuntimeContainer, GPUCount: 1}, models.RuntimeService, "GPU"},
		{"unknown", &Service{cluster: clusterOn{}}, &models.Application{RuntimeKind: models.RuntimeContainer}, "vm", "unknown runtime"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := tc.s.PreflightRuntime(tc.app, RuntimeChange{Runtime: tc.to})
			if !hasMsg(p.Blockers, tc.want) {
				t.Errorf("blockers = %v, want one containing %q", p.Blockers, tc.want)
			}
		})
	}
}

// A blocked change is refused before anything is saved.
func TestChangeRuntimeRefusesBlocked(t *testing.T) {
	s := &Service{cluster: clusterOff{}}
	app := &models.Application{RuntimeKind: models.RuntimeContainer}
	_, err := s.ChangeRuntime(app, RuntimeChange{Runtime: models.RuntimeService}, nil)
	var rc *RuntimeChangeError
	if err == nil || !errors.As(err, &rc) || rc.Code() != "RUNTIME_CHANGE_BLOCKED" {
		t.Fatalf("err = %v, want a RuntimeChangeError", err)
	}
	if app.RuntimeKind != models.RuntimeContainer {
		t.Error("a refused change must leave the app alone")
	}
}
