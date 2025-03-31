package style

import (
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	Title         lipgloss.Style
	Info          lipgloss.Style
	ContentHeader lipgloss.Style
	Success       lipgloss.Style
	Error         lipgloss.Style
	Tag           lipgloss.Style
	Doc           lipgloss.Style
	PromptList    lipgloss.Style
	Viewport      lipgloss.Style
	InputLabel    lipgloss.Style
	InputView     lipgloss.Style
	Footer        lipgloss.Style
	StatusLine    lipgloss.Style
}

func New() Styles {
	return Styles{
		Title:         lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true).MarginLeft(2),
		Info:          lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA")).Background(lipgloss.Color("#4D7EA8")).Padding(0, 1),
		ContentHeader: lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA")).Background(lipgloss.Color("#2D9862")).Bold(true).Padding(0, 1),
		Success:       lipgloss.NewStyle().Foreground(lipgloss.Color("#2D9862")).Bold(true),
		Error:         lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F87")).Bold(true),
		Tag:           lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Background(lipgloss.Color("#242424")).Padding(0, 1).MarginRight(1),
		Doc:           lipgloss.NewStyle().Margin(1, 2),
		PromptList:    lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#4D7EA8")).Padding(1, 0),
		Viewport:      lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#4D7EA8")).Padding(1, 0),
		InputLabel:    lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true),
		InputView:     lipgloss.NewStyle().Padding(0, 1),
		Footer:        lipgloss.NewStyle().MarginTop(1),
		StatusLine:    lipgloss.NewStyle().Height(1),
	}
}
