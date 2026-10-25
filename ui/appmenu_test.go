package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestAppmenuCloseSubmenuWithNoItems(t *testing.T) {
	// The menu is empty whenever Tailscale is not running, and the left key
	// still calls CloseSubmenu.
	var appmenu Appmenu
	appmenu.SetItems([]*AppmenuItem{})

	appmenu.CloseSubmenu()

	if appmenu.IsSubmenuOpen() {
		t.Fatal("expected submenu to be closed")
	}
}

func TestAppmenuRenderChildSubmenuPreviewFitsWidth(t *testing.T) {
	// Network Devices shows the selected device's details next to the device
	// list as soon as the submenu is open, which doesn't fit beside the main
	// menu in an 80 column terminal.
	details := &Submenu{}
	details.SetItems([]SubmenuItem{
		&TitleSubmenuItem{Label: "Machine Details"},
		&LabeledSubmenuItem{Label: "Tailscale IPv4", AdditionalLabel: "100.64.0.2"},
	})

	networkDevices := &AppmenuItem{Label: "Network Devices"}
	networkDevices.Submenu.SetItems([]SubmenuItem{
		&LabeledSubmenuItem{Label: "nas", AdditionalLabel: "Linux", ChildSubmenu: details},
	})

	var appmenu Appmenu
	appmenu.SetItems([]*AppmenuItem{networkDevices})
	appmenu.Activate()

	const width = 80
	rendered := appmenu.Render(20, width)
	for _, line := range strings.Split(rendered, "\n") {
		if lineWidth := lipgloss.Width(line); lineWidth > width {
			t.Fatalf("rendered line is %d columns wide, want at most %d: %q", lineWidth, width, line)
		}
	}
	if !strings.Contains(rendered, "Machine Details") {
		t.Fatal("expected the machine details to be shown")
	}
}
