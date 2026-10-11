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
	"testing"

	"github.com/miabi-io/miabi/internal/models"
)

func newAddressService() (*Service, *memStore) {
	store := &memStore{
		def: models.Cluster{ID: 1, IsDefault: true, Name: "default"},
		others: map[uint]models.Cluster{
			2: {ID: 2, Name: "frankfurt", Mode: models.ClusterModeStandalone, ManagerServerID: 21, IngressServerID: 21},
			3: {ID: 3, Name: "paris", Mode: models.ClusterModeSwarm, ManagerServerID: 31, IngressServerID: 31, IngressIP: "192.0.2.3"},
		},
	}
	nodes := &fakeNodes{servers: []models.Server{
		{ID: 1, IsLocal: true, ClusterID: 1},
		{ID: 21, ClusterID: 2},
		{ID: 31, ClusterID: 3},
		{ID: 32, ClusterID: 3},
	}}
	return &Service{store: store, nodes: nodes}, store
}

func TestEveryClusterHasAnEditablePublicAddress(t *testing.T) {
	s, store := newAddressService()

	ip, host := " 203.0.113.9 ", "Edge.Example.com."
	if _, err := s.UpdateCluster(models.DefaultClusterID, ClusterPatch{IngressIP: &ip, IngressHostname: &host}); err != nil {
		t.Fatalf("update: %v", err)
	}
	if gotIP, gotHost := s.IngressAddress(models.DefaultClusterID); gotIP != "203.0.113.9" || gotHost != "edge.example.com" {
		t.Errorf("default cluster address = %q, %q; want the normalized IP and hostname", gotIP, gotHost)
	}
	if store.def.IngressIP != "203.0.113.9" {
		t.Errorf("stored IP = %q", store.def.IngressIP)
	}

	bad := "not-an-ip"
	if _, err := s.UpdateCluster(2, ClusterPatch{IngressIP: &bad}); !errors.Is(err, ErrInvalidIngressIP) {
		t.Errorf("bad IP err = %v, want ErrInvalidIngressIP", err)
	}
	machine := "orbstack"
	if _, err := s.UpdateCluster(2, ClusterPatch{IngressHostname: &machine}); !errors.Is(err, ErrInvalidIngressHostname) {
		t.Errorf("machine hostname err = %v, want ErrInvalidIngressHostname: a CNAME needs a DNS name", err)
	}
}

func TestAClusterLearnsItsAddressOnlyFromItsIngressNode(t *testing.T) {
	s, store := newAddressService()

	s.LearnIngressIP(21, "10.0.0.5:4711")
	if store.others[2].IngressIP != "" {
		t.Errorf("a private source address was learned: %q", store.others[2].IngressIP)
	}
	s.LearnIngressIP(21, "203.0.113.21:4711")
	if got := store.others[2].IngressIP; got != "203.0.113.21" {
		t.Errorf("standalone cluster address = %q, want its node's public source IP", got)
	}
	s.LearnIngressIP(21, "203.0.113.99")
	if got := store.others[2].IngressIP; got != "203.0.113.21" {
		t.Errorf("a learned address replaced one already set: %q", got)
	}
	s.LearnIngressIP(32, "203.0.113.32")
	if got := store.others[3].IngressIP; got != "192.0.2.3" {
		t.Errorf("a swarm member that is not the ingress node set the address: %q", got)
	}
	s.LearnIngressIP(1, "203.0.113.1")
	if store.def.IngressIP != "" {
		t.Errorf("the control-plane node's source IP was learned: %q", store.def.IngressIP)
	}
}
