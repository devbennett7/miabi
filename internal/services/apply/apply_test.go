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

	d "github.com/miabi-io/miabi/internal/declarative"
)

func app(name string, env map[string]string) d.Resource {
	return d.Resource{Kind: d.KindApplication, Metadata: d.Meta{Name: name}, Application: &d.ApplicationSpec{Env: env}}
}

func db(name string) d.Resource {
	return d.Resource{Kind: d.KindDatabase, Metadata: d.Meta{Name: name}, Database: &d.DatabaseSpec{Engine: "postgres"}}
}

func TestRawConsumerApp(t *testing.T) {
	tests := []struct {
		name string
		all  []d.Resource
		dep  string
		want string
	}{
		{
			name: "app that references the database wins",
			all: []d.Resource{
				db("db"),
				app("web", map[string]string{"DATABASE_URL": "{{ .databases.db.uri }}"}),
				app("worker", nil),
			},
			dep:  "db",
			want: "web",
		},
		{
			name: "sole app is the consumer even without a reference",
			all:  []d.Resource{db("db"), app("web", nil)},
			dep:  "db",
			want: "web",
		},
		{
			name: "ambiguous: multiple apps, none reference it",
			all:  []d.Resource{db("db"), app("web", nil), app("worker", nil)},
			dep:  "db",
			want: "",
		},
		{
			name: "the referencing app is picked among several",
			all: []d.Resource{
				db("cache"),
				app("web", nil),
				app("worker", map[string]string{"REDIS_URL": "{{ .databases.cache.uri }}"}),
			},
			dep:  "cache",
			want: "worker",
		},
		{
			name: "no reference and no app yields no consumer",
			all:  []d.Resource{db("db")},
			dep:  "db",
			want: "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := rawConsumerApp(tc.all, tc.dep); got != tc.want {
				t.Errorf("rawConsumerApp(%q) = %q, want %q", tc.dep, got, tc.want)
			}
		})
	}
}
