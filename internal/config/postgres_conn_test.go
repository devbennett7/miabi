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

package config

import "testing"

func TestPostgresConnDiscreteFields(t *testing.T) {
	d := DatabaseConfig{host: "db", user: "miabi", password: "secret", name: "miabi", port: 5432, sslMode: "require"}
	host, port, name, user, pass, ssl := d.PostgresConn()
	if host != "db" || port != 5432 || name != "miabi" || user != "miabi" || pass != "secret" || ssl != "require" {
		t.Fatalf("discrete fields not resolved: %s %d %s %s %s %s", host, port, name, user, pass, ssl)
	}
}

func TestPostgresConnURLOverrides(t *testing.T) {
	d := DatabaseConfig{
		host: "fallback", port: 5432, name: "fallback", user: "fallback", password: "fallback", sslMode: "disable",
		url: "postgres://u:p@pghost:6543/appdb?sslmode=require",
	}
	host, port, name, user, pass, ssl := d.PostgresConn()
	if host != "pghost" || port != 6543 || name != "appdb" || user != "u" || pass != "p" || ssl != "require" {
		t.Fatalf("url not parsed: %s %d %s %s %s %s", host, port, name, user, pass, ssl)
	}
}

func TestPostgresConnBadURLFallsBack(t *testing.T) {
	d := DatabaseConfig{host: "db", port: 5432, name: "miabi", user: "miabi", password: "x", sslMode: "disable", url: "::not-a-url::"}
	host, port, _, _, _, _ := d.PostgresConn()
	if host != "db" || port != 5432 {
		t.Fatalf("bad url should fall back to discrete fields, got host=%s port=%d", host, port)
	}
}
