package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/neuralinkcorp/tsui/ui"
)

// A line in the help screen.
type helpEntry struct {
	// Keys as reported by tea.KeyMsg.String().
	keys        []string
	description string
}

// Every key handled in Update. help_test.go checks that each of these actually does something.
var helpEntries = []helpEntry{
	{[]string{"up", "k", "w"}, "Move up"},
	{[]string{"down", "j", "s"}, "Move down"},
	{[]string{"right", "l", "d"}, "Open a menu or item details"},
	{[]string{"left", "h", "a"}, "Go back"},
	{[]string{"enter", " "}, "Open a menu or select an item"},
	{[]string{"."}, "Connect, disconnect, or log in"},
	{[]string{"esc"}, "Go back, or quit from the main menu"},
	{[]string{"?"}, "Show or hide this help"},
	{[]string{"q", "ctrl+c"}, "Quit"},
}

// Display names for keys that don't read well as-is.
var helpKeyNames = map[string]string{
	"up":    "↑",
	"down":  "↓",
	"left":  "←",
	"right": "→",
	" ":     "space",
}

// Render the help screen, centered in and clipped to the given size.
func renderHelp(width int, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}

	labels := make([]string, len(helpEntries))
	labelWidth := 0
	for i, entry := range helpEntries {
		names := make([]string, len(entry.keys))
		for j, key := range entry.keys {
			names[j] = key
			if name, ok := helpKeyNames[key]; ok {
				names[j] = name
			}
		}
		labels[i] = strings.Join(names, ", ")
		labelWidth = max(labelWidth, lipgloss.Width(labels[i]))
	}

	keyStyle := lipgloss.NewStyle().
		Foreground(ui.Primary).
		Bold(true).
		Width(labelWidth + 3)

	lines := []string{
		lipgloss.NewStyle().
			Bold(true).
			Render("Keyboard Shortcuts") +
			lipgloss.NewStyle().
				Faint(true).
				Render(" (? or esc to close)"),
		"",
	}
	for i, entry := range helpEntries {
		lines = append(lines, keyStyle.Render(labels[i])+entry.description)
	}

	// JoinVertical pads the lines to the same width so the columns stay aligned when centered.
	return lipgloss.NewStyle().
		MaxWidth(width).
		MaxHeight(height).
		Render(lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center,
			lipgloss.JoinVertical(lipgloss.Left, lines...)))
}
