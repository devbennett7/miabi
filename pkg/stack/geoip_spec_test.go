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

import "testing"

// The GeoIP database is bound into the gateway (at Goma's GOMA_GEOIP_DB default
// path) only when it has been provisioned — binding a missing host path would
// make Docker create a directory there.
func TestGatewaySpecBindsGeoIPOnlyWhenProvisioned(t *testing.T) {
	m := Defaults("miabi/miabi:1.4.0")

	bind := func() (source string, readOnly, found bool) {
		for _, b := range gatewaySpec(m, ContainerGateway, m.Images.Gateway).Binds {
			if b.Target == gomaGeoIPFile {
				return b.Source, b.ReadOnly, true
			}
		}
		return "", false, false
	}

	if _, _, found := bind(); found {
		t.Fatal("GeoIP database bound when none was provisioned")
	}

	const host = "/etc/miabi/country.mmdb"
	m.gatewayHostGeoIP = host
	source, readOnly, found := bind()
	if !found {
		t.Fatal("GeoIP database not bound after provisioning")
	}
	if source != host || !readOnly {
		t.Fatalf("GeoIP bind wrong: source=%q readOnly=%v, want %q read-only", source, readOnly, host)
	}
}
