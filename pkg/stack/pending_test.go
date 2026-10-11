// Copyright (c) 2026 Jonas Kaninda.
//
// This file is part of Miabi.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package stack

import (
	"context"
	"slices"
	"testing"

	"github.com/miabi-io/miabi/pkg/stack/docker"
)

type specEngine struct {
	docker.Client
	running map[string]docker.Container
}

func (e *specEngine) InspectContainer(_ context.Context, name string) (docker.Container, error) {
	c, ok := e.running[name]
	if !ok {
		return docker.Container{}, docker.ErrNotFound
	}
	return c, nil
}

// convergedEngine reports every component running with the spec hash Converge would have stamped.
func convergedEngine(s *Service, m *Manifest) *specEngine {
	e := &specEngine{running: map[string]docker.Container{}}
	for _, c := range s.components(m) {
		spec := c.Build(m, c.Name, *c.Image(m))
		e.running[c.Name] = docker.Container{State: "running", Labels: map[string]string{docker.LabelSpecHash: specHash(spec)}}
	}
	return e
}

// Several env edits must add up to one recreate of the component they touch, and nothing else.
func TestPendingNamesOnlyTheComponentsAnEditChanged(t *testing.T) {
	m := Defaults("miabi/miabi:1.10.0")
	m.Domain = "miabi.example.com"
	if err := m.Normalize(); err != nil {
		t.Fatal(err)
	}
	s := &Service{log: func(string, ...any) {}}
	e := convergedEngine(s, m)
	s.dc = e

	if got, err := s.Pending(context.Background(), m); err != nil || len(got) != 0 {
		t.Fatalf("converged stack: Pending = %v, %v; want nothing", got, err)
	}

	m.Env["MIABI_SMTP_HOST"] = "smtp.example.com"
	m.Env["MIABI_SMTP_PORT"] = "587"
	got, err := s.Pending(context.Background(), m)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got, []string{ContainerControlPlane}) {
		t.Errorf("after two control-plane env edits Pending = %v, want only %s", got, ContainerControlPlane)
	}

	delete(e.running, ContainerGateway)
	if got, _ := s.Pending(context.Background(), m); !slices.Contains(got, ContainerGateway) {
		t.Errorf("a missing gateway is not pending: %v", got)
	}
}
