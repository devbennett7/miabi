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

package placement

import (
	"errors"
	"strings"
	"testing"
)

func TestReach(t *testing.T) {
	swarm := func(id uint) bool { return id == 2 }
	label := func(id uint) string { return map[uint]string{1: "eu-central", 2: "eu-east"}[id] }
	app := func(cluster, server uint) Site {
		return Site{Kind: "application", Name: "api", ClusterID: cluster, ServerID: server}
	}
	db := func(cluster, server uint) Site {
		return Site{Kind: "database", Name: "pg", ClusterID: cluster, ServerID: server}
	}

	if err := Reach(app(1, 5), db(1, 5), swarm, label); err != nil {
		t.Errorf("same node: %v", err)
	}
	if err := Reach(app(2, 5), db(2, 6), swarm, label); err != nil {
		t.Errorf("two nodes of a swarm: %v", err)
	}
	if err := Reach(app(1, 5), db(1, 6), swarm, label); !errors.Is(err, ErrCrossNode) {
		t.Errorf("two nodes without a swarm = %v", err)
	}
	err := Reach(app(2, 5), db(1, 5), swarm, label)
	if !errors.Is(err, ErrCrossLocation) {
		t.Fatalf("two locations = %v", err)
	}
	if want := "application api is in eu-east, database pg is in eu-central"; !strings.HasPrefix(err.Error(), want) {
		t.Errorf("message = %q", err)
	}
}
