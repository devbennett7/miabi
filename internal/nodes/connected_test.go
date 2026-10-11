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

package nodes

import (
	"testing"
	"time"
)

func TestConnectedSince(t *testing.T) {
	c := NewClients(1, &stubClient{})

	if at, ok := c.ConnectedSince(1); !ok || !at.IsZero() {
		t.Fatalf("local ConnectedSince = (%v, %v); want (zero, true)", at, ok)
	}
	if _, ok := c.ConnectedSince(2); ok {
		t.Fatal("a node with no client reports connected")
	}

	before := time.Now()
	c.SetRemote(2, &stubClient{})
	if at, ok := c.ConnectedSince(2); !ok || at.Before(before) {
		t.Fatalf("ConnectedSince after connect = (%v, %v); want a time from the connect", at, ok)
	}

	c.RemoveRemote(2)
	if _, ok := c.ConnectedSince(2); ok {
		t.Fatal("a disconnected node still reports connected")
	}
}
