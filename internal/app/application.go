package app

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/renatogalera/promptgen/internal/config"
	"github.com/renatogalera/promptgen/internal/domain/prompt"
	"github.com/renatogalera/promptgen/internal/domain/xml"
	"github.com/renatogalera/promptgen/internal/storage/yaml"
	"github.com/renatogalera/promptgen/internal/ui/keymap"
	"github.com/renatogalera/promptgen/internal/ui/style"
	"github.com/renatogalera/promptgen/pkg/clipboard"
	"github.com/renatogalera/promptgen/pkg/filter"
)

type Application struct {
	model Model
}

func NewApplication(promptFile string) *Application {
	promptService := prompt.NewService()
	yamlRepo := yaml.NewRepository(promptFile, promptService)
	xmlFormatter := xml.NewFormatter()
	clipboardManager := clipboard.New()
	s := spinner.New(
		spinner.WithSpinner(spinner.Dot),
		spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("205"))),
	)
	delegate := list.NewDefaultDelegate()
	l := list.New(nil, delegate, 0, 0)
	l.Title = "AI Prompts"
	l.Styles.Title = style.New().Title

	l.SetFilteringEnabled(true)
	l.SetShowFilter(true)
	l.SetShowStatusBar(true)
	l.SetShowHelp(false)
	l.SetShowPagination(true)
	l.SetShowTitle(true)

	l.Filter = filter.MultiTokenSubstringFilter

	vp := viewport.New(0, 0)
	vp.Style = style.New().Viewport

	h := help.New()
	h.ShowAll = true

	app := &Application{
		model: Model{
			keyMap:        keymap.New(),
			help:          h,
			list:          l,
			spinner:       s,
			viewport:      vp,
			promptFile:    promptFile,
			state:         config.StateLoading,
			showHelp:      true,
			variables:     make(map[string]string),
			styles:        style.New(),
			promptService: promptService,
			yamlRepo:      yamlRepo,
			xmlFormatter:  xmlFormatter,
			clipboardMgr:  clipboardManager,
		},
	}

	return app
}

func (a *Application) Init() tea.Cmd {
	return a.model.Init()
}

func (a *Application) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	modelUpdated, cmd := a.model.Update(msg)

	if updatedModel, ok := modelUpdated.(Model); ok {
		a.model = updatedModel
	}

	return a, cmd
}

func (a *Application) View() string {
	return a.model.View()
}
