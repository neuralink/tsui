package main

import (
	"strings"
	"testing"
	"time"

	"tailscale.com/ipn/ipnstate"
)

func TestMachineDetailsSubmenuExpiredKey(t *testing.T) {
	keyExpiry := time.Now().Add(-72 * time.Hour)
	peer := &ipnstate.PeerStatus{
		DNSName:   "laptop.example.ts.net.",
		KeyExpiry: &keyExpiry,
		Expired:   true,
	}

	rendered := buildMachineDetailsSubmenu(peer).Render(true, 40, 60)

	for _, line := range strings.Split(rendered, "\n") {
		if strings.Contains(line, "Key Expiry") {
			if !strings.Contains(line, "Expired") {
				t.Fatalf("expected expired key to be shown as Expired, got %q", line)
			}
			return
		}
	}
	t.Fatalf("Key Expiry item not rendered:\n%s", rendered)
}
