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

package controlmanager

import (
	"testing"

	"github.com/miabi-io/miabi/internal/models"
	"github.com/miabi-io/miabi/internal/services/settings"
)

func TestEffectivePolicy(t *testing.T) {
	cases := []struct {
		name     string
		policy   models.ReconcilePolicy
		platform Mode
		want     Mode
	}{
		{"inherit follows the platform", models.ReconcileInherit, ModeEnforce, ModeEnforce},
		{"empty policy follows the platform", "", ModeObserve, ModeObserve},
		{"off opts out of an enforcing platform", models.ReconcileOff, ModeEnforce, ModeOff},
		{"observe reports without acting", models.ReconcileObserve, ModeEnforce, ModeObserve},
		{"enforce opts in while the platform observes", models.ReconcileEnforce, ModeObserve, ModeEnforce},
		// A platform that is off runs no sweep at all, so an app cannot enforce its way back in.
		{"enforce cannot override a platform that is off", models.ReconcileEnforce, ModeOff, ModeOff},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := effectivePolicy(&models.Application{ReconcilePolicy: c.policy}, c.platform)
			if got != c.want {
				t.Fatalf("effectivePolicy(%q, %q) = %q; want %q", c.policy, c.platform, got, c.want)
			}
		})
	}
}

// An app exempted by its own policy is not watched at all: no finding, no event, nothing to act on.
func TestAppWithPolicyOffIsNotWatched(t *testing.T) {
	h := newHarness()
	h.enforcing()
	app := containerApp(7, 1)
	app.ReconcilePolicy = models.ReconcileOff
	*h.apps = fakeApps{app}
	h.releases[7] = "c7"

	if st := h.confirm(t); len(st.Findings) != 0 || len(h.events.events) != 0 {
		t.Fatalf("findings=%+v events=%v; an exempted app must not be watched", st.Findings, h.events.events)
	}
}

// Observe on one app holds while the rest of the platform enforces: reported, never redeployed.
func TestAppWithPolicyObserveIsReportedButNotRedeployed(t *testing.T) {
	h := newHarness()
	e := h.enforcing()
	app := containerApp(7, 1)
	app.ReconcilePolicy = models.ReconcileObserve
	*h.apps = fakeApps{app}
	h.releases[7] = "c7"

	st := h.confirm(t)
	if len(st.Findings) != 1 {
		t.Fatalf("findings = %+v; want it still reported", st.Findings)
	}
	if len(e.redeploy.calls) != 0 {
		t.Fatalf("redeploys = %+v; want none for an observe-only app", e.redeploy.calls)
	}
}

// And the other direction: one app may be brought back while the platform only observes.
func TestAppWithPolicyEnforceActsWhileThePlatformObserves(t *testing.T) {
	h := newHarness()
	e := h.enforcing()
	h.settings[settings.KeyControlManagerMode] = string(ModeObserve)
	app := containerApp(7, 1)
	app.ReconcilePolicy = models.ReconcileEnforce
	*h.apps = fakeApps{app}
	h.releases[7] = "c7"

	h.confirm(t)
	if len(e.redeploy.calls) != 1 || e.redeploy.calls[0].appID != 7 {
		t.Fatalf("redeploys = %+v; want app 7 brought back on its own policy", e.redeploy.calls)
	}
}
