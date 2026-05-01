package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewPromptModel_DefaultValue(t *testing.T) {
	m := newPromptModel("label", "placeholder", "default")
	if m.input.Value() != "default" {
		t.Errorf("expected default value, got %q", m.input.Value())
	}
}

func TestPromptModel_Init(t *testing.T) {
	m := newPromptModel("label", "placeholder", "")
	cmd := m.Init()
	if cmd == nil {
		t.Error("expected non-nil init command (blink)")
	}
}

func TestPromptModel_View(t *testing.T) {
	m := newPromptModel("My Label", "ph", "")
	view := m.View()
	if !strings.Contains(view, "My Label") {
		t.Errorf("expected label in view, got: %s", view)
	}
}

func TestPromptModel_Update_Enter_UsesValue(t *testing.T) {
	m := newPromptModel("label", "placeholder", "")
	m.input.SetValue("typed")

	result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	pm := result.(promptModel)

	if !pm.submitted {
		t.Error("expected submitted = true after Enter")
	}
	if pm.value != "typed" {
		t.Errorf("expected value 'typed', got %q", pm.value)
	}
}

func TestPromptModel_Update_Enter_FallsBackToPlaceholder(t *testing.T) {
	m := newPromptModel("label", "myplaceholder", "")

	result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	pm := result.(promptModel)

	if pm.value != "myplaceholder" {
		t.Errorf("expected placeholder value, got %q", pm.value)
	}
}

func TestPromptModel_Update_Escape(t *testing.T) {
	m := newPromptModel("label", "ph", "")

	result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	pm := result.(promptModel)

	if pm.submitted {
		t.Error("expected submitted = false after Escape")
	}
}

func TestPromptModel_Update_CtrlC(t *testing.T) {
	m := newPromptModel("label", "ph", "")

	result, _ := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	pm := result.(promptModel)

	if pm.submitted {
		t.Error("expected submitted = false after Ctrl+C")
	}
}
