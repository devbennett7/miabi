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

import "testing"

func TestClusterSlotsCapEachClusterSeparately(t *testing.T) {
	slots := newClusterSlots(2)

	release, ok1 := slots.acquire(1)
	_, ok2 := slots.acquire(1)
	if !ok1 || !ok2 {
		t.Fatal("the first two deploys of a cluster were refused")
	}
	if _, ok := slots.acquire(1); ok {
		t.Fatal("a third concurrent deploy was let into a cluster capped at two")
	}
	if _, ok := slots.acquire(2); !ok {
		t.Fatal("a busy cluster blocked a deploy to another cluster")
	}

	release()
	release()
	if _, ok := slots.acquire(1); !ok {
		t.Fatal("a released slot could not be reused")
	}
	if _, ok := slots.acquire(1); ok {
		t.Fatal("releasing the same slot twice freed two")
	}
}
