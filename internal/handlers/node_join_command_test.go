// SPDX-FileCopyrightText: 2026 Jonas Kaninda
// SPDX-License-Identifier: AGPL-3.0-or-later

package handlers

import "testing"

func TestAgentScriptCommand(t *testing.T) {
	got := agentScriptCommand("https://miabi.example.com", "mbn_xxxxxxxx")
	want := "curl -fsSL https://get.miabi.io/agent | \\\n" +
		"  MIABI_CONTROL_URL=https://miabi.example.com MIABI_NODE_TOKEN=mbn_xxxxxxxx bash"
	if got != want {
		t.Fatalf("script command =\n%s\nwant\n%s", got, want)
	}
}
