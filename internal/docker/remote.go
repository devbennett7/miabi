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

package docker

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/moby/moby/client"
)

// DialFunc opens a new connection to a remote Docker engine. For multi-node, it
// returns a fresh stream over the node's agent tunnel; each Docker API request
// gets its own stream so streaming endpoints (logs/stats/events/build) work.
type DialFunc func(ctx context.Context) (net.Conn, error)

// NewRemote builds a Docker client that talks to a remote engine over an arbitrary transport (the
// agent tunnel), reusing engineClient so every service works against a remote node unchanged. The
// dummy "http://" host makes the SDK speak plain HTTP; the dialer returns a tunnel stream.
func NewRemote(dial DialFunc) (Client, error) {
	httpc := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return dial(ctx)
			},
			// Docker's streaming/hijacked endpoints hold a connection open; don't
			// cap the per-request lifetime here (callers pass their own ctx).
			MaxIdleConns:          10,
			IdleConnTimeout:       90 * time.Second,
			ResponseHeaderTimeout: 0,
			DisableCompression:    true,
		},
	}
	// Order matters: WithHost runs sockets.ConfigureTransport on the *current* client's transport (which would
	// clobber our DialContext). Apply it first against the throwaway default client, then install our http
	// client last so our tunnel dialer survives.
	cli, err := client.New(
		client.WithHost("http://docker.invalid"),
		client.WithHTTPClient(httpc),
	)
	if err != nil {
		return nil, err
	}
	return &engineClient{cli: cli}, nil
}
