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

import (
	"slices"
	"testing"
)

const countryHeaderEnv = "GOMA_GEOIP_COUNTRY_HEADER=true"

func TestCountryHeaderReachesTheGateway(t *testing.T) {
	m := testManifest()
	if err := m.Normalize(); err != nil {
		t.Fatal(err)
	}
	if slices.Contains(gatewayEnv(m), countryHeaderEnv) {
		t.Fatal("country header sent while gateway.countryHeader is off")
	}

	m.Install.GatewayCountryHeader = true
	if err := m.Normalize(); err != nil {
		t.Fatal(err)
	}
	if env := gatewayEnv(m); !slices.Contains(env, countryHeaderEnv) {
		t.Errorf("gateway env lacks %s (got %v)", countryHeaderEnv, env)
	}
}

func TestTogglingCountryHeaderRecreatesTheGateway(t *testing.T) {
	m := testManifest()
	before := specHash(gatewaySpec(m, ContainerGateway, m.Images.Gateway))
	m.Install.GatewayCountryHeader = true
	if specHash(gatewaySpec(m, ContainerGateway, m.Images.Gateway)) == before {
		t.Error("the gateway spec hash ignored countryHeader, so converge would leave the old gateway running")
	}
}

func TestASpecCountryHeaderIsRefusedInGatewayEnv(t *testing.T) {
	m := testManifest()
	m.Install.GatewayCountryHeader = true
	m.Gateway.Env["GOMA_GEOIP_COUNTRY_HEADER"] = "false"
	if err := m.Normalize(); err == nil {
		t.Error("gateway.env was allowed to contradict gateway.countryHeader")
	}
}
