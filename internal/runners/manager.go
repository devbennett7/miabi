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

// Package runners maintains the live tunnels of connected build/pipeline runners. It mirrors the
// node-agent connection manager but is thinner: a runner uses its OWN Docker/BuildKit daemon, so
// the tunnel exists only to lease jobs to it and stream logs and status back.
package runners

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hashicorp/yamux"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/wstunnel"
	"github.com/miabi-io/miabi/internal/models"
)

// heartbeatInterval refreshes a connected runner's last-seen, matching the node
// agent's cadence.
const heartbeatInterval = 30 * time.Second

// connectionState is what the runner service needs to reflect a live tunnel.
// Implemented by runner.Service (MarkConnected/MarkDisconnected); an interface
// keeps this package free of a service import cycle.
type connectionState interface {
	MarkConnected(id uint, f models.RunnerFacts)
	MarkDisconnected(id uint)
}

// Manager accepts runner tunnels and tracks each runner's live session and
// online/offline status.
type Manager struct {
	state connectionState

	mu          sync.Mutex
	sessions    map[uint]*yamux.Session
	statusHooks []func(id uint, online bool)
}

// NewManager creates a runner connection manager backed by the runner service
// (for status persistence).
func NewManager(state connectionState) *Manager {
	return &Manager{state: state, sessions: map[uint]*yamux.Session{}}
}

// Handle owns an authenticated runner WebSocket: it builds the tunnel, marks the runner online,
// and blocks until the tunnel closes so the HTTP handler keeps the connection open. The caller has
// already validated the registration token.
func (m *Manager) Handle(r *models.Runner, f models.RunnerFacts, ws *websocket.Conn) {
	id := r.ID
	// Control plane OPENS streams to the runner (to dispatch job leases), so the
	// runner side ACCEPTS them — mirrors the node/agent role split.
	sess, err := wstunnel.Client(ws)
	if err != nil {
		logger.Error("failed to start runner tunnel", "runner", id, "error", err)
		_ = ws.Close()
		return
	}

	m.serve(r, f, sess)
}

// serve tracks one runner session from connect until its tunnel closes.
func (m *Manager) serve(r *models.Runner, f models.RunnerFacts, sess *yamux.Session) {
	id := r.ID
	m.replace(id, sess)
	m.state.MarkConnected(id, f)
	m.notify(id, true)
	logger.Info("runner connected", "runner", id, "name", r.Name, "version", f.Version, "ip", f.RemoteIP)

	stop := make(chan struct{})
	go m.heartbeat(id, f, stop)

	<-sess.CloseChan() // blocks until the tunnel dies (keepalive timeout or close)
	close(stop)

	// A tunnel superseded by a reconnect closes after its replacement is live, and must not mark
	// the runner offline under it.
	if m.forget(id, sess) {
		m.state.MarkDisconnected(id)
		m.notify(id, false)
	}
	logger.Info("runner disconnected", "runner", id, "name", r.Name)
}

// AddStatusListener registers a hook fired when a runner's tunnel comes up or goes down, so the
// console can show a newly enrolled runner without a reload. Hooks run in order on the tunnel's
// goroutine, keeping a quick drop and reconnect in sequence, so they must not block.
func (m *Manager) AddStatusListener(fn func(id uint, online bool)) {
	m.mu.Lock()
	m.statusHooks = append(m.statusHooks, fn)
	m.mu.Unlock()
}

func (m *Manager) notify(id uint, online bool) {
	m.mu.Lock()
	hooks := append([]func(uint, bool){}, m.statusHooks...)
	m.mu.Unlock()
	for _, fn := range hooks {
		fn(id, online)
	}
}

// Connected reports whether a runner currently has a live tunnel.
func (m *Manager) Connected(id uint) bool {
	m.mu.Lock()
	_, ok := m.sessions[id]
	m.mu.Unlock()
	return ok
}

// Session returns a connected runner's multiplexed session (for the P3 lease
// dispatcher to open streams). ok is false when the runner is offline.
func (m *Manager) Session(id uint) (*yamux.Session, bool) {
	m.mu.Lock()
	sess, ok := m.sessions[id]
	m.mu.Unlock()
	return sess, ok
}

// Disconnect tears down a runner's tunnel (e.g. when it is deleted or disabled).
func (m *Manager) Disconnect(id uint) {
	m.mu.Lock()
	sess := m.sessions[id]
	m.mu.Unlock()
	if sess != nil {
		_ = sess.Close()
	}
}

func (m *Manager) heartbeat(id uint, f models.RunnerFacts, stop <-chan struct{}) {
	t := time.NewTicker(heartbeatInterval)
	defer t.Stop()
	for {
		select {
		case <-stop:
			return
		case <-t.C:
			m.state.MarkConnected(id, f)
		}
	}
}

// replace registers the new session, closing any prior one (a reconnect
// supersedes a stale tunnel).
func (m *Manager) replace(id uint, sess *yamux.Session) {
	m.mu.Lock()
	old := m.sessions[id]
	m.sessions[id] = sess
	m.mu.Unlock()
	if old != nil {
		_ = old.Close()
	}
}

// forget drops sess and reports whether it was still the runner's live session.
func (m *Manager) forget(id uint, sess *yamux.Session) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.sessions[id] != sess {
		return false
	}
	delete(m.sessions, id)
	return true
}
