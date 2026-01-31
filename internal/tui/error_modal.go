package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ErrorModal displays error messages in a modal overlay
type ErrorModal struct {
	err     error
	visible bool
	width   int
	height  int
	styles  *Styles
}

// NewErrorModal creates a new ErrorModal
func NewErrorModal(styles *Styles) ErrorModal {
	return ErrorModal{
		styles: styles,
	}
}

// Init initializes the error modal
func (m ErrorModal) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m ErrorModal) Update(msg tea.Msg) (ErrorModal, tea.Cmd) {
	switch msg := msg.(type) {
	case ErrorMsg:
		m.err = msg.Err
		m.visible = true

	case tea.KeyMsg:
		if m.visible {
			// Any key dismisses the modal
			m.visible = false
			m.err = nil
		}

	case WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View renders the error modal
func (m ErrorModal) View() string {
	if !m.visible || m.err == nil {
		return ""
	}

	// Build the modal content
	title := m.styles.ErrorTitle.Render("Error")
	message := m.err.Error()

	// Wrap message if too long
	maxWidth := m.width - 10
	if maxWidth < 20 {
		maxWidth = 20
	}
	message = wrapText(message, maxWidth)

	// Build modal content
	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		m.styles.ErrorMessage.Render(message),
		"",
		m.styles.StatusValue.Render("Press any key to dismiss"),
	)

	// Apply border
	modal := m.styles.ErrorBorder.Render(content)

	// Center the modal
	modalWidth := StringWidth(modal)
	modalHeight := strings.Count(modal, "\n") + 1

	// Calculate position
	x := (m.width - modalWidth) / 2
	y := (m.height - modalHeight) / 2

	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	// Position the modal (using newlines and spaces)
	var sb strings.Builder
	for range y {
		sb.WriteString("\n")
	}
	for range x {
		sb.WriteString(" ")
	}
	sb.WriteString(modal)

	return sb.String()
}

// IsVisible returns whether the modal is visible
func (m *ErrorModal) IsVisible() bool {
	return m.visible
}

// Show shows the modal with an error
func (m *ErrorModal) Show(err error) {
	m.err = err
	m.visible = true
}

// Hide hides the modal
func (m *ErrorModal) Hide() {
	m.visible = false
	m.err = nil
}

// SetSize sets the modal's available size
func (m *ErrorModal) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// wrapText wraps text to fit within maxWidth
func wrapText(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return text
	}

	var result strings.Builder
	var lineWidth int

	words := strings.Fields(text)
	for i, word := range words {
		wordWidth := StringWidth(word)

		if lineWidth+wordWidth+1 > maxWidth && lineWidth > 0 {
			result.WriteString("\n")
			lineWidth = 0
		} else if i > 0 && lineWidth > 0 {
			result.WriteString(" ")
			lineWidth++
		}

		result.WriteString(word)
		lineWidth += wordWidth
	}

	return result.String()
}
