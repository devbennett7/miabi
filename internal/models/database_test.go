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

package models

import (
	"reflect"
	"testing"
)

func TestDatabaseInstanceNetworkNames(t *testing.T) {
	const fallback = "miabi"

	tests := []struct {
		name string
		inst DatabaseInstance
		want []string
	}{
		{
			name: "legacy: no primary, no attachments falls back to gateway",
			inst: DatabaseInstance{},
			want: []string{fallback},
		},
		{
			name: "primary only",
			inst: DatabaseInstance{NetworkName: "mb-ws1-abc"},
			want: []string{"mb-ws1-abc"},
		},
		{
			name: "primary listed first, then extra attached networks",
			inst: DatabaseInstance{
				NetworkName: "mb-ws1-abc",
				Networks:    []Network{{DockerName: "mb-ws1-abc"}, {DockerName: "mb-ws1-custom"}},
			},
			want: []string{"mb-ws1-abc", "mb-ws1-custom"},
		},
		{
			name: "attached network reachable even when not the pinned primary",
			inst: DatabaseInstance{
				NetworkName: "miabi",
				Networks:    []Network{{DockerName: "mb-ws1-abc"}},
			},
			want: []string{"miabi", "mb-ws1-abc"},
		},
		{
			name: "blank docker names are skipped",
			inst: DatabaseInstance{
				NetworkName: "mb-ws1-abc",
				Networks:    []Network{{DockerName: ""}, {DockerName: "mb-ws1-abc"}},
			},
			want: []string{"mb-ws1-abc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.inst.NetworkNames(fallback)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("NetworkNames() = %v, want %v", got, tt.want)
			}
		})
	}
}
