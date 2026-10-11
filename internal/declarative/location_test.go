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

func stackSet(location string) *d.ResourceSet {
	set := d.NewResourceSet()
	set.Add(d.Resource{
		APIVersion: d.APIVersion, Kind: d.KindStack,
		Metadata: d.Meta{Name: "shop"}, Stack: &d.StackSpec{Placement: &d.PlacementSpec{Location: location}},
	})
	return set
}

func stackChange(desired, actual string) (d.Change, bool) {
	for _, c := range d.BuildPlan(stackSet(desired), stackSet(actual), d.PlanOptions{}).Changes {
		if c.Kind == d.KindStack && c.Name == "shop" {
			return c, true
		}
	}
	return d.Change{}, false
}

// A manifest that names no location must not drift against a resource that has one.
func TestUnstatedLocationIsNotDrift(t *testing.T) {
	if c, ok := stackChange("", "eu-central"); ok && c.Action != d.ActionNoop {
		t.Errorf("change = %+v, want no drift", c)
	}
}

func TestChangedLocationPlansAnUpdate(t *testing.T) {
	c, _ := stackChange("eu-east", "eu-central")
	if c.Action != d.ActionUpdate || len(c.Fields) != 1 || c.Fields[0].Field != "placement.location" {
		t.Errorf("change = %+v, want an update of the location alone", c)
	}
}

func TestLocationParses(t *testing.T) {
	set, err := d.Parse([]byte("apiVersion: miabi.io/v1\nkind: Volume\nmetadata:\n  name: data\nspec:\n  placement:\n    location: eu-east\n"))
	if err != nil {
		t.Fatal(err)
	}
	r, _ := set.Get("Volume/data")
	if r.Volume == nil || r.Volume.Location() != "eu-east" {
		t.Errorf("volume = %+v, want location eu-east", r.Volume)
	}
}
