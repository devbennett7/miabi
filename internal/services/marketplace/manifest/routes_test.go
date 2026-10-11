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

package manifest

import (
	"strings"
	"testing"
)

func routesTemplate(routes string) string {
	return `
apiVersion: miabi.io/v1
kind: Template
metadata: {name: a, displayName: A, version: "1.0.0"}
inputs:
  - key: public_url
    label: URL
    type: string
applications:
  - name: web
    image: nginx
    ports:
      - container: 8080
        scheme: http
  - name: worker
    image: nginx
routes:
` + routes
}

func TestRoutesValidation(t *testing.T) {
	cases := map[string]string{
		"":                `  - app: web` + "\n" + `    host: "{{ .inputs.public_url }}"`,
		"unknown app":     `  - app: nope` + "\n" + `    host: "{{ .inputs.public_url }}"`,
		"no port":         `  - app: worker` + "\n" + `    host: "{{ .inputs.public_url }}"`,
		"undeclared port": `  - app: web` + "\n" + `    host: "{{ .inputs.public_url }}"` + "\n" + `    port: 9999`,
		"literal host":    `  - app: web` + "\n" + `    host: status.example.com`,
	}
	for want, routes := range cases {
		_, err := Parse([]byte(routesTemplate(routes + "\n")))
		if want == "" {
			if err != nil {
				t.Errorf("valid route rejected: %v", err)
			}
			continue
		}
		if err == nil {
			t.Errorf("%s: accepted", want)
		}
	}
}

func TestRouteHost(t *testing.T) {
	cases := map[string]string{
		"https://Status.Example.com/":       "status.example.com",
		"https://status.example.com:8443/x": "status.example.com",
		"status.example.com":                "status.example.com",
		"status.example.com.":               "status.example.com",
		"":                                  "",
		"http://localhost:3000":             "",
		"http://10.0.0.5":                   "",
		"https://[::1]:443":                 "",
		"intranet":                          "",
	}
	for in, want := range cases {
		if got := RouteHost(in); got != want {
			t.Errorf("RouteHost(%q) = %q, want %q", in, got, want)
		}
	}
	if strings.Contains(RouteHost("https://a.example.com/path?q=1"), "/") {
		t.Error("path leaked into the host")
	}
}
