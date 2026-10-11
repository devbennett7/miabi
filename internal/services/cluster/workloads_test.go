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
	"errors"
	"strings"
	"testing"
)

// A refusal has to say what is in the way and what to do instead, not just that the node is busy.
func TestNodeWorkloadsErrorNamesTheWayForward(t *testing.T) {
	for _, tt := range []struct {
		err  *NodeWorkloadsError
		want []string
	}{
		{&NodeWorkloadsError{Node: "lyon", Action: EnableSwarmAction, Apps: 2, Databases: 1, Volumes: 3},
			[]string{`"lyon"`, "2 apps, 1 database and 3 volumes", "Swarm cannot be enabled", "default cluster", "empty node"}},
		{&NodeWorkloadsError{Node: "lyon", Action: JoinAction, Databases: 1},
			[]string{"runs 1 database, so it cannot join", "default cluster"}},
		{&NodeWorkloadsError{Node: "lyon", Action: LeaveAction, Apps: 1, Volumes: 1},
			[]string{"1 app and 1 volume", "cannot leave"}},
	} {
		if !errors.Is(tt.err, ErrNodeHasWorkloads) {
			t.Errorf("%v does not match ErrNodeHasWorkloads", tt.err)
		}
		for _, want := range tt.want {
			if !strings.Contains(tt.err.Error(), want) {
				t.Errorf("message is missing %q: %s", want, tt.err)
			}
		}
	}
}
