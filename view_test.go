package main

import (
	"strings"
	"testing"

	"github.com/neuralinkcorp/tsui/libts"
	"tailscale.com/ipn"
)

func TestViewNeedsMachineAuth(t *testing.T) {
	m := model{
		terminalWidth:  100,
		terminalHeight: 30,
		state: libts.State{
			BackendState: ipn.NeedsMachineAuth,
			Prefs:        &ipn.Prefs{},
		},
	}

	view := m.View()

	for _, want := range []string{
		"Device Approval Required",
		"https://login.tailscale.com/admin/machines",
		"continue automatically",
	} {
		if !strings.Contains(view, want) {
			t.Errorf("view is missing %q:\n%s", want, view)
		}
	}
}
