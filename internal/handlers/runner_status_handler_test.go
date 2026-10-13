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

package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jkaninda/okapi"
	"github.com/miabi-io/miabi/internal/models"
	"github.com/miabi-io/miabi/internal/services/eventbus"
)

// syncRecorder lets the test read the stream while the handler is still writing it.
type syncRecorder struct {
	mu  sync.Mutex
	rec *httptest.ResponseRecorder
}

func (s *syncRecorder) Header() http.Header { return s.rec.Header() }
func (s *syncRecorder) WriteHeader(code int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rec.WriteHeader(code)
}
func (s *syncRecorder) Write(b []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rec.Write(b)
}
func (s *syncRecorder) Flush() {}
func (s *syncRecorder) body() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rec.Body.String()
}

func TestRunnerStatusStreamFiltersToViewer(t *testing.T) {
	bus := eventbus.New()
	h := &RunnerHandler{bus: bus}
	ws, other := uint(5), uint(6)

	ctx, cancel := context.WithCancel(context.Background())
	w := &syncRecorder{rec: httptest.NewRecorder()}
	c := okapi.NewContext(okapi.New(), w, httptest.NewRequest(http.MethodGet, "/events", nil).WithContext(ctx))

	listed := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		done <- h.streamStatus(c,
			func() ([]models.Runner, error) {
				close(listed)
				return []models.Runner{{ID: 1, Name: "snap-runner", WorkspaceID: &ws}}, nil
			},
			func(r *models.Runner) bool { return r.WorkspaceID != nil && *r.WorkspaceID == ws },
		)
	}()
	<-listed

	publish := func(r models.Runner) {
		bus.Publish(runnerStatusTopic, eventbus.Event{Type: "status", Data: r})
	}
	publish(models.Runner{ID: 2, Name: "other-ws-runner", WorkspaceID: &other, Connected: true})
	publish(models.Runner{ID: 3, Name: "shared-runner", Connected: true})
	publish(models.Runner{ID: 4, Name: "own-runner", WorkspaceID: &ws, Connected: true})

	deadline := time.Now().Add(2 * time.Second)
	for !strings.Contains(w.body(), "own-runner") {
		if time.Now().After(deadline) {
			t.Fatalf("own runner's event never streamed; body:\n%s", w.body())
		}
		time.Sleep(10 * time.Millisecond)
	}
	cancel()
	if err := <-done; err != nil {
		t.Fatalf("streamStatus: %v", err)
	}

	body := w.body()
	if !strings.Contains(body, `"type":"snapshot"`) || !strings.Contains(body, "snap-runner") {
		t.Errorf("stream did not open with the snapshot:\n%s", body)
	}
	if strings.Index(body, "snap-runner") > strings.Index(body, "own-runner") {
		t.Errorf("snapshot was not sent first:\n%s", body)
	}
	for _, leaked := range []string{"other-ws-runner", "shared-runner"} {
		if strings.Contains(body, leaked) {
			t.Errorf("stream leaked %s to another viewer:\n%s", leaked, body)
		}
	}
}
