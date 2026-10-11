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

package platformbackup

import (
	"context"
	"testing"
	"time"
)

// A recovery point must outlive the request that asked for it. Tenant capture used to run inline in the
// HTTP handler's context, so the moment that request ended every artifact still in flight died with
// "context canceled". CreateSet and RetrySet now detach, which is what this pins.
func TestWithoutCancelSurvivesTheRequest(t *testing.T) {
	req, cancel := context.WithCancel(context.Background())
	work := context.WithoutCancel(req)

	cancel() // the HTTP handler returned

	if req.Err() == nil {
		t.Fatal("the request context should be cancelled")
	}
	if err := work.Err(); err != nil {
		t.Fatalf("the detached context was cancelled with the request: %v", err)
	}

	// And it stays usable for the kind of wait a helper container needs.
	select {
	case <-work.Done():
		t.Fatal("the detached context reported done")
	case <-time.After(10 * time.Millisecond):
	}
}

// A deadline on the request must not travel either: an artifact that takes longer
// than the HTTP timeout is normal, not a failure.
func TestWithoutCancelDropsTheDeadline(t *testing.T) {
	req, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	work := context.WithoutCancel(req)

	time.Sleep(20 * time.Millisecond)
	if req.Err() == nil {
		t.Fatal("the request context should have expired")
	}
	if _, ok := work.Deadline(); ok {
		t.Fatal("the detached context inherited the request's deadline")
	}
	if err := work.Err(); err != nil {
		t.Fatalf("the detached context expired: %v", err)
	}
}
