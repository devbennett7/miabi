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

package config_test

import (
	"testing"

	"github.com/miabi-io/miabi/internal/config"
	"github.com/miabi-io/miabi/pkg/stack"
)

func TestInstallerAndControlPlaneAgreeOnLogLevels(t *testing.T) {
	for _, lvl := range []string{"debug", "info", "warn", "warning", "error"} {
		m := stack.Defaults("miabi/miabi:1.4.0")
		m.Domain = "miabi.example.com"
		m.Env = map[string]string{"MIABI_LOG_LEVEL": lvl}
		if err := m.Normalize(); err != nil {
			t.Errorf("the installer rejects %q: %v", lvl, err)
			continue
		}
		if _, err := (&config.Config{LogLevel: lvl}).LogLevelFor(); err != nil {
			t.Errorf("the installer accepts %q but the control plane rejects it (%v) — "+
				"the install would succeed and then crash-loop", lvl, err)
		}
	}

	for _, lvl := range []string{"off", "none", "silent", "verbose", "trace"} {
		if _, err := (&config.Config{LogLevel: lvl}).LogLevelFor(); err == nil {
			continue
		}
		m := stack.Defaults("miabi/miabi:1.4.0")
		m.Domain = "miabi.example.com"
		m.Env = map[string]string{"MIABI_LOG_LEVEL": lvl}
		if err := m.Normalize(); err == nil {
			t.Errorf("the installer accepts %q but the control plane rejects it — "+
				"the install would report success and the panel would never come up", lvl)
		}
	}
}
