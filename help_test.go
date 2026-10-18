package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/neuralinkcorp/tsui/libts"
	"github.com/neuralinkcorp/tsui/ui"
	"tailscale.com/ipn"
)

// Build the KeyMsg whose String() is key.
func keyMsg(key string) tea.KeyMsg {
	switch key {
	case "up":
		return tea.KeyMsg{Type: tea.KeyUp}
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}
	case "left":
		return tea.KeyMsg{Type: tea.KeyLeft}
	case "right":
		return tea.KeyMsg{Type: tea.KeyRight}
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEscape}
	case "ctrl+c":
		return tea.KeyMsg{Type: tea.KeyCtrlC}
	case " ":
		return tea.KeyMsg{Type: tea.KeySpace}
	}
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
}

func press(m model, key string) (model, tea.Cmd) {
	next, cmd := m.Update(keyMsg(key))
	return next.(model), cmd
}

func isQuit(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	_, ok := cmd().(tea.QuitMsg)
	return ok
}

// A model with a three item menu, with the cursor on the middle item.
func newKeyTestModel() model {
	items := []*ui.AppmenuItem{{Label: "One"}, {Label: "Two"}, {Label: "Three"}}
	for _, item := range items {
		item.Submenu.SetItems([]ui.SubmenuItem{
			&ui.TitleSubmenuItem{Label: "submenu " + item.Label},
			&ui.LabeledSubmenuItem{Label: "item"},
		})
	}

	m := model{
		terminalWidth:  80,
		terminalHeight: 24,
		state:          libts.State{BackendState: ipn.Stopped},
	}
	m.menu.SetItems(items)
	m.menu.CursorDown()
	return m
}

// The submenu preview shown next to the main menu belongs to the selected item.
func selectedItem(m model) string {
	for _, label := range []string{"One", "Two", "Three"} {
		if strings.Contains(m.menu.Render(10, 80), "submenu "+label) {
			return label
		}
	}
	return ""
}

func TestHelpKeysAreHandled(t *testing.T) {
	moves := func(want string) func(t *testing.T, key string) {
		return func(t *testing.T, key string) {
			m, _ := press(newKeyTestModel(), key)
			if got := selectedItem(m); got != want {
				t.Errorf("selected %q, want %q", got, want)
			}
		}
	}
	opens := func(t *testing.T, key string) {
		m, _ := press(newKeyTestModel(), key)
		if !m.menu.IsSubmenuOpen() {
			t.Error("expected submenu to open")
		}
	}
	closes := func(t *testing.T, key string) {
		m := newKeyTestModel()
		m.menu.Activate()
		m, _ = press(m, key)
		if m.menu.IsSubmenuOpen() {
			t.Error("expected submenu to close")
		}
	}
	quits := func(t *testing.T, key string) {
		if _, cmd := press(newKeyTestModel(), key); !isQuit(cmd) {
			t.Error("expected quit")
		}
	}

	checks := map[string]func(t *testing.T, key string){
		"up":     moves("One"),
		"k":      moves("One"),
		"w":      moves("One"),
		"down":   moves("Three"),
		"j":      moves("Three"),
		"s":      moves("Three"),
		"right":  opens,
		"l":      opens,
		"d":      opens,
		"enter":  opens,
		" ":      opens,
		"left":   closes,
		"h":      closes,
		"a":      closes,
		"q":      quits,
		"ctrl+c": quits,
		"esc": func(t *testing.T, key string) {
			closes(t, key)
			quits(t, key)
		},
		".": func(t *testing.T, key string) {
			// Don't run the command, it would talk to the daemon.
			if _, cmd := press(newKeyTestModel(), key); cmd == nil {
				t.Error("expected a command to bring Tailscale up")
			}
		},
		"?": func(t *testing.T, key string) {
			if m, _ := press(newKeyTestModel(), key); !m.showHelp {
				t.Error("expected help to open")
			}
		},
	}

	for _, entry := range helpEntries {
		for _, key := range entry.keys {
			if got := keyMsg(key).String(); got != key {
				t.Fatalf("keyMsg(%q).String() = %q", key, got)
			}
			check, ok := checks[key]
			if !ok {
				t.Errorf("help lists %q but this test doesn't check it", key)
				continue
			}
			t.Run(key, func(t *testing.T) { check(t, key) })
		}
	}
}

func TestHelpToggle(t *testing.T) {
	m := newKeyTestModel()

	m, _ = press(m, "?")
	if !m.showHelp || !strings.Contains(m.View(), "Keyboard Shortcuts") {
		t.Fatal("expected ? to open help")
	}

	// Menu keys are ignored while help is open.
	m, _ = press(m, "down")
	m, _ = press(m, "enter")
	if selectedItem(m) != "Two" || m.menu.IsSubmenuOpen() {
		t.Error("menu changed while help was open")
	}

	m, _ = press(m, "?")
	if m.showHelp {
		t.Error("expected ? to close help")
	}

	m, _ = press(m, "?")
	m, cmd := press(m, "esc")
	if m.showHelp || isQuit(cmd) {
		t.Error("expected esc to close help without quitting")
	}

	m, _ = press(m, "?")
	if _, cmd := press(m, "q"); !isQuit(cmd) {
		t.Error("expected q to quit while help is open")
	}
}

func TestHelpRenderSizes(t *testing.T) {
	sizes := [][2]int{{0, 0}, {1, 1}, {5, 2}, {20, 5}, {40, 10}, {80, 24}, {200, 60}}
	for _, size := range sizes {
		width, height := size[0], size[1]

		m := newKeyTestModel()
		m.terminalWidth, m.terminalHeight = width, height
		m.showHelp = true
		_ = m.View()

		help := renderHelp(width, height)
		if help == "" {
			continue
		}
		if lipgloss.Height(help) > height {
			t.Errorf("%dx%d: help is %d lines tall", width, height, lipgloss.Height(help))
		}
		for _, line := range strings.Split(help, "\n") {
			if lipgloss.Width(line) > width {
				t.Errorf("%dx%d: help line is %d wide", width, height, lipgloss.Width(line))
			}
		}
	}
}
