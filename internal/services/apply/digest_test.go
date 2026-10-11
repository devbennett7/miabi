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

package apply

import (
	"testing"

	"github.com/miabi-io/miabi/internal/models"
)

func TestRunningDigest(t *testing.T) {
	const (
		pin = "sha256:aaaa"
		old = "sha256:bbbb"
	)
	pinned := "registry.example.com/ws_3/web@" + pin
	cases := []struct {
		name   string
		latest *models.Deployment
		active *models.Release
		want   string
	}{
		{
			"converged",
			&models.Deployment{Status: models.DeploymentSucceeded, Image: pinned},
			&models.Release{Image: pinned},
			pin,
		},
		{
			"failed deploy leaves the old release running",
			&models.Deployment{Status: models.DeploymentFailed, Image: pinned},
			&models.Release{Image: "registry.example.com/ws_3/web:latest"},
			pin + " (not running; running an unpinned image)",
		},
		{
			"failed deploy over an older pin",
			&models.Deployment{Status: models.DeploymentFailed, Image: pinned},
			&models.Release{Image: "registry.example.com/ws_3/web@" + old},
			pin + " (not running; running " + old + ")",
		},
		{
			"deploy to the pin in flight counts as converged",
			&models.Deployment{Status: models.DeploymentDeploying, Image: pinned},
			&models.Release{Image: "registry.example.com/ws_3/web:latest"},
			pin,
		},
		{
			"deploy in flight to something else does not",
			&models.Deployment{Status: models.DeploymentPending, Image: "registry.example.com/ws_3/web:latest"},
			&models.Release{Image: "registry.example.com/ws_3/web@" + old},
			pin + " (not running; running " + old + ")",
		},
		{
			"release digest wins over the image ref",
			&models.Deployment{Status: models.DeploymentSucceeded, Image: "registry.example.com/ws_3/web:latest"},
			&models.Release{Image: "registry.example.com/ws_3/web:latest", Digest: pin},
			pin,
		},
		{"never deployed", nil, nil, pin + " (not running; running an unpinned image)"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := runningDigest(pin, c.latest, c.active); got != c.want {
				t.Errorf("runningDigest = %q, want %q", got, c.want)
			}
		})
	}
}
