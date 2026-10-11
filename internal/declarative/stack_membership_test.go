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

package declarative_test

import (
	"testing"

	d "github.com/miabi-io/miabi/internal/declarative"
)

func stackBundle(appStack string) string {
	body := `
apiVersion: miabi.io/v1
kind: Stack
metadata: { name: shop }
---
apiVersion: miabi.io/v1
kind: Application
metadata: { name: web }
spec:
  image: ghcr.io/org/web
`
	if appStack != "" {
		body += "  stack: " + appStack + "\n"
	}
	return body
}

func stackDiff(t *testing.T, desired, actual string) (d.FieldDiff, bool) {
	t.Helper()
	ds, err := d.Parse([]byte(desired))
	if err != nil {
		t.Fatal(err)
	}
	as, err := d.Parse([]byte(actual))
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range d.BuildPlan(ds, as, d.PlanOptions{}).Changes {
		for _, f := range c.Fields {
			if f.Field == "stack" {
				return f, true
			}
		}
	}
	return d.FieldDiff{}, false
}

// Joining a stack from a manifest has to plan as a change of the app.
func TestStatedStackConverges(t *testing.T) {
	f, ok := stackDiff(t, stackBundle("shop"), stackBundle(""))
	if !ok || f.From != "" || f.To != "shop" {
		t.Errorf("stack diff = %+v (planned %v), want \"\" -> shop", f, ok)
	}
}

// A manifest silent about the stack leaves a membership set in the console alone.
func TestOmittedStackIsNotDrift(t *testing.T) {
	if f, ok := stackDiff(t, stackBundle(""), stackBundle("shop")); ok {
		t.Errorf("a manifest without a stack planned %+v", f)
	}
}
