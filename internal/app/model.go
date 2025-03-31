package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
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
)

type Model struct {
	keyMap         keymap.KeyMap
	help           help.Model
	list           list.Model
	spinner        spinner.Model
	viewport       viewport.Model
	textInputs     []textinput.Model
	styles         style.Styles
	promptService  *prompt.Service
	yamlRepo       *yaml.Repository
	xmlFormatter   *xml.Formatter
	clipboardMgr   *clipboard.Manager
	inputLabels    []string
	prompts        prompt.PromptCollection
	selectedPrompt prompt.Prompt
	state          config.AppState
	promptFile     string
	statusMessage  string
	statusCmd      tea.Cmd
	showHelp       bool
	activeInput    int
	variables      map[string]string
	width          int
	height         int
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type promptsLoadedMsg struct {
	prompts prompt.PromptCollection
	items   []list.Item
}

type copyDoneMsg struct{}
type promptSavedMsg struct{}
type statusMsg struct{ message string }
type clearStatusMsg struct{}

func (m Model) Init() tea.Cmd {

	return tea.Batch(
		m.spinner.Tick,
		m.loadPromptsCmd(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

		m = m.updateWindowSize(msg)

	case tea.KeyMsg:

		switch {
		case keymap.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit

		case keymap.Matches(msg, m.keyMap.Help):
			m.showHelp = !m.showHelp

			m = m.updateWindowSize(tea.WindowSizeMsg{Width: m.width, Height: m.height})
		}

		switch m.state {
		case config.StatePromptList:

			if m.list.FilterState() == list.Filtering {
				m.list, cmd = m.list.Update(msg)
				cmds = append(cmds, cmd)
			} else {

				switch {
				case keymap.Matches(msg, m.keyMap.Select):
					selectedListItem := m.list.SelectedItem()
					if selectedListItem != nil {
						if item, ok := selectedListItem.(prompt.Item); ok {
							m.selectedPrompt = item.Prompt
							m.state = config.StatePromptView
							m.viewport.SetContent(m.selectedPrompt.Content)
							m.viewport.GotoTop()
							m.help.ShowAll = false
						}
					}
				case keymap.Matches(msg, m.keyMap.Create):
					m.state = config.StatePromptCreation
					m.inputLabels = []string{"Title", "Tags (comma-sep)", "Description", "Content"}
					m.initTextInputs(len(m.inputLabels))
					m.activeInput = 0
					if len(m.textInputs) > 0 {

						cmds = append(cmds, m.textInputs[m.activeInput].Focus())
					}
					m.help.ShowAll = true
				default:

					m.list, cmd = m.list.Update(msg)
					cmds = append(cmds, cmd)
				}
			}
		case config.StatePromptView:

			switch {
			case keymap.Matches(msg, m.keyMap.Back):
				m.state = config.StatePromptList
				m.selectedPrompt = prompt.Prompt{}
				m.viewport.SetContent("")
				m.help.ShowAll = true
			case keymap.Matches(msg, m.keyMap.Copy):

				cmd = m.copyToClipboardCmd(m.selectedPrompt)
				cmds = append(cmds, cmd)
			case keymap.Matches(msg, m.keyMap.Select) && len(m.selectedPrompt.Variables) > 0:
				m.state = config.StateVariableInput
				m.inputLabels = m.selectedPrompt.Variables
				m.initTextInputs(len(m.inputLabels))
				m.activeInput = 0
				if len(m.textInputs) > 0 {
					cmds = append(cmds, m.textInputs[m.activeInput].Focus())
				}
				m.help.ShowAll = true
			default:

				m.viewport, cmd = m.viewport.Update(msg)
				cmds = append(cmds, cmd)
			}
		case config.StateVariableInput, config.StatePromptCreation:

			switch {
			case keymap.Matches(msg, m.keyMap.Cancel):
				if len(m.textInputs) > 0 && m.activeInput < len(m.textInputs) {
					m.textInputs[m.activeInput].Blur()
				}

				if m.state == config.StatePromptCreation {
					m.state = config.StatePromptList
					m.help.ShowAll = true
				} else {
					m.state = config.StatePromptView
					m.help.ShowAll = false
				}

				m.textInputs = nil
				m.inputLabels = nil
				m.variables = make(map[string]string)
			case keymap.Matches(msg, m.keyMap.Confirm):
				if m.state == config.StateVariableInput {

					cmd = m.handleVariableInputConfirm()
					cmds = append(cmds, cmd)
				} else {

					cmds = append(cmds, m.saveNewPromptCmd())
				}

			case keymap.Matches(msg, m.keyMap.Up),
				keymap.Matches(msg, m.keyMap.Down),
				keymap.MatchesKeyType(msg, tea.KeyTab),
				keymap.MatchesKeyType(msg, tea.KeyShiftTab):

				if len(m.textInputs) > 1 {
					cmd = m.handleInputFocus(msg)
					cmds = append(cmds, cmd)
				} else if len(m.textInputs) == 1 {

					m.textInputs[m.activeInput], cmd = m.textInputs[m.activeInput].Update(msg)
					cmds = append(cmds, cmd)
				}
			default:

				if len(m.textInputs) > 0 && m.activeInput < len(m.textInputs) {
					m.textInputs[m.activeInput], cmd = m.textInputs[m.activeInput].Update(msg)
					cmds = append(cmds, cmd)
				}
			}
		}

	case errMsg:
		m.statusMessage = m.styles.Error.Render(fmt.Sprintf("ERROR: %v", msg.err))
		m.statusCmd = m.clearStatusCmd()
		cmds = append(cmds, m.statusCmd)

	case promptsLoadedMsg:
		m.state = config.StatePromptList
		m.prompts = msg.prompts
		m.list.SetItems(msg.items)
		m.help.ShowAll = true

	case copyDoneMsg:
		m.statusMessage = m.styles.Success.Render("Copied to clipboard!")
		m.statusCmd = m.clearStatusCmd()
		cmds = append(cmds, m.statusCmd)

	case promptSavedMsg:
		m.statusMessage = m.styles.Success.Render("Prompt saved!")

		m.state = config.StateLoading
		if len(m.textInputs) > 0 && m.activeInput < len(m.textInputs) {
			m.textInputs[m.activeInput].Blur()
		}
		m.textInputs = nil
		m.inputLabels = nil
		m.help.ShowAll = true

		cmds = append(cmds, m.loadPromptsCmd(), m.clearStatusCmd())

	case statusMsg:

		m.statusMessage = msg.message
		m.statusCmd = m.clearStatusCmd()
		cmds = append(cmds, m.statusCmd)

	case clearStatusMsg:

		m.statusMessage = ""
		m.statusCmd = nil

	case spinner.TickMsg:

		if m.state == config.StateLoading {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	default:
		switch m.state {
		case config.StatePromptView:
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		case config.StatePromptList:
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		case config.StateVariableInput, config.StatePromptCreation:

			for i := range m.textInputs {
				m.textInputs[i], cmd = m.textInputs[i].Update(msg)
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var s strings.Builder

	s.WriteString(m.styles.Title.Render("PromptGen") + "\n\n")

	switch m.state {
	case config.StateLoading:

		loadingStyle := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).PaddingTop(m.height/2 - 1)
		s.WriteString(loadingStyle.Render(fmt.Sprintf("%s Loading prompts...", m.spinner.View())))
	case config.StatePromptList:
		s.WriteString(m.styles.PromptList.Render(m.list.View()))
	case config.StatePromptView:
		s.WriteString(m.renderPromptViewHeader())
		s.WriteString(m.styles.Viewport.Render(m.viewport.View()))
	case config.StateVariableInput, config.StatePromptCreation:
		s.WriteString(m.renderInputForm())
	default:
		s.WriteString(m.styles.Error.Render("Unknown application state!"))
	}

	statusLine := m.renderStatusLine()
	helpView := m.renderHelpView()

	footerContent := lipgloss.JoinVertical(lipgloss.Left, statusLine, helpView)

	contentHeight := lipgloss.Height(s.String())
	footerHeight := lipgloss.Height(footerContent)

	paddingHeight := m.height - contentHeight - footerHeight - 1
	if paddingHeight < 0 {
		paddingHeight = 0
	}

	s.WriteString(strings.Repeat("\n", paddingHeight))
	s.WriteString(m.styles.Footer.Render(footerContent))

	return s.String()
}

func (m Model) updateWindowSize(msg tea.WindowSizeMsg) Model {
	m.width = msg.Width
	m.height = msg.Height
	m.help.Width = msg.Width

	titleHeight := lipgloss.Height(m.styles.Title.Render("PromptGen") + "\n\n")
	statusHeight := lipgloss.Height(m.renderStatusLine())
	helpHeight := lipgloss.Height(m.renderHelpView())
	footerHeight := statusHeight + helpHeight + lipgloss.Height(m.styles.Footer.Render(""))

	availableHeight := m.height - titleHeight - footerHeight

	listStyle := m.styles.PromptList
	listVPadding := listStyle.GetVerticalPadding()
	listHPadding := listStyle.GetHorizontalPadding()
	m.list.SetSize(m.width-listHPadding, availableHeight-listVPadding)

	viewHeaderHeight := 0
	if m.state == config.StatePromptView {
		viewHeaderHeight = lipgloss.Height(m.renderPromptViewHeader())
	}
	viewportStyle := m.styles.Viewport
	vpVPadding := viewportStyle.GetVerticalPadding()
	vpHPadding := viewportStyle.GetHorizontalPadding()

	m.viewport.Width = m.width - vpHPadding
	m.viewport.Height = availableHeight - viewHeaderHeight - vpVPadding

	if len(m.textInputs) > 0 {

		formStyle := m.styles.Doc
		formHPadding := formStyle.GetHorizontalPadding()
		inputPromptWidth := lipgloss.Width(m.textInputs[0].Prompt)

		inputWidth := m.width - formHPadding - inputPromptWidth - 2
		minInputWidth := 20
		if inputWidth < minInputWidth {
			inputWidth = minInputWidth
		}
		for i := range m.textInputs {
			m.textInputs[i].Width = inputWidth
		}
	}

	return m
}

func (m Model) renderStatusLine() string {

	return m.styles.StatusLine.Render(m.statusMessage)
}

func (m Model) renderHelpView() string {
	if m.showHelp {
		return m.help.View(m.keyMap)
	}

	return m.help.ShortHelpView(m.keyMap.ShortHelp())
}

func (m Model) renderPromptViewHeader() string {
	var header strings.Builder
	header.WriteString(m.styles.Title.Render(m.selectedPrompt.Title) + "\n")
	if len(m.selectedPrompt.Tags) > 0 {
		var styledTags []string
		for _, t := range m.selectedPrompt.Tags {
			styledTags = append(styledTags, m.styles.Tag.Render(t))
		}
		header.WriteString(strings.Join(styledTags, " ") + "\n")
	}
	if m.selectedPrompt.Description != "" {
		header.WriteString(m.styles.Info.Render(m.selectedPrompt.Description) + "\n")
	}

	header.WriteString("\n" + m.styles.ContentHeader.Render("Content:") + "\n")
	return header.String()
}

func (m Model) renderInputForm() string {
	var form strings.Builder
	viewTitle := "New Prompt"
	if m.state == config.StateVariableInput {
		viewTitle = fmt.Sprintf("Variables for: %s", m.selectedPrompt.Title)
	}
	form.WriteString(m.styles.Title.Render(viewTitle) + "\n\n")

	for i := range m.textInputs {
		label := ""
		if i < len(m.inputLabels) {
			label = m.inputLabels[i]
		}
		form.WriteString(m.styles.InputLabel.Render(label+":") + "\n")

		form.WriteString(m.styles.InputView.Render(m.textInputs[i].View()) + "\n\n")
	}

	form.WriteString(
		lipgloss.NewStyle().Faint(true).Render(
			"Use Tab/Shift+Tab or Up/Down to navigate, Enter to confirm, Esc to cancel.",
		),
	)

	return m.styles.Doc.Render(form.String())
}

func (m *Model) initTextInputs(count int) {
	m.textInputs = make([]textinput.Model, count)

	formStyle := m.styles.Doc
	formHPadding := formStyle.GetHorizontalPadding()
	inputPromptWidth := lipgloss.Width("┃ ")
	inputWidth := m.width - formHPadding - inputPromptWidth - 2
	minInputWidth := 20
	if inputWidth < minInputWidth {
		inputWidth = minInputWidth
	}

	for i := range m.textInputs {
		t := textinput.New()
		t.Width = inputWidth
		t.Prompt = "┃ "
		t.PromptStyle = m.styles.InputLabel
		t.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
		if i == 0 {
			t.Focus()
		}
		m.textInputs[i] = t
	}
	m.activeInput = 0
}

func (m *Model) handleInputFocus(msg tea.KeyMsg) tea.Cmd {
	current := m.activeInput
	numInputs := len(m.textInputs)
	if numInputs <= 1 {
		return nil
	}

	var focusCmds []tea.Cmd

	if current >= 0 && current < numInputs {
		m.textInputs[current].Blur()
	}

	switch {
	case keymap.MatchesKeyType(msg, tea.KeyTab), keymap.Matches(msg, m.keyMap.Down):
		m.activeInput = (current + 1) % numInputs
	case keymap.MatchesKeyType(msg, tea.KeyShiftTab), keymap.Matches(msg, m.keyMap.Up):
		m.activeInput = (current - 1 + numInputs) % numInputs
	}

	if m.activeInput >= 0 && m.activeInput < numInputs {
		focusCmds = append(focusCmds, m.textInputs[m.activeInput].Focus())
	}

	return tea.Batch(focusCmds...)
}

func (m *Model) handleVariableInputConfirm() tea.Cmd {

	m.variables = make(map[string]string)
	for i, label := range m.inputLabels {
		if i < len(m.textInputs) {
			m.variables[label] = m.textInputs[i].Value()
		}
	}

	tempPrompt := m.selectedPrompt
	tempPrompt.Content = m.promptService.ReplaceVariables(tempPrompt.Content, m.variables)

	m.state = config.StatePromptView
	m.help.ShowAll = false

	if len(m.textInputs) > 0 && m.activeInput < len(m.textInputs) {
		m.textInputs[m.activeInput].Blur()
	}
	m.textInputs = nil
	m.inputLabels = nil

	return m.copyToClipboardCmd(tempPrompt)
}

func (m *Model) copyToClipboardCmd(p prompt.Prompt) tea.Cmd {
	return func() tea.Msg {

		xmlOutput, err := m.xmlFormatter.FormatWithDoc(p)
		if err != nil {

			return statusMsg{message: m.styles.Error.Render(fmt.Sprintf("XML Generation error: %v", err))}
		}

		if err := m.clipboardMgr.Copy(xmlOutput); err != nil {

			return statusMsg{message: m.styles.Error.Render(fmt.Sprintf("Clipboard error: %v", err))}
		}

		return copyDoneMsg{}
	}
}

func (m *Model) saveNewPromptCmd() tea.Cmd {
	return func() tea.Msg {

		title := ""
		tagString := ""
		description := ""
		content := ""
		if len(m.textInputs) > 0 {
			title = m.textInputs[0].Value()
		}
		if len(m.textInputs) > 1 {
			tagString = m.textInputs[1].Value()
		}
		if len(m.textInputs) > 2 {
			description = m.textInputs[2].Value()
		}
		if len(m.textInputs) > 3 {
			content = m.textInputs[3].Value()
		}

		if strings.TrimSpace(title) == "" {
			return statusMsg{message: m.styles.Error.Render("Title cannot be empty")}
		}
		if strings.TrimSpace(content) == "" {
			return statusMsg{message: m.styles.Error.Render("Content cannot be empty")}
		}

		newPrompt := m.promptService.CreatePrompt(
			title,
			tagString,
			description,
			content,
		)

		if err := m.yamlRepo.SavePrompt(newPrompt); err != nil {
			return statusMsg{message: m.styles.Error.Render(fmt.Sprintf("Error saving prompt: %v", err))}
		}

		return promptSavedMsg{}
	}
}

func (m *Model) loadPromptsCmd() tea.Cmd {
	return func() tea.Msg {

		collection, err := m.yamlRepo.LoadPrompts()
		if err != nil {

			return errMsg{err: fmt.Errorf("failed to load prompts from %s: %w", m.promptFile, err)}
		}

		items := collection.ToItems()

		return promptsLoadedMsg{prompts: collection, items: items}
	}
}

func (m *Model) clearStatusCmd() tea.Cmd {

	return tea.Tick(config.StatusMessageTimeout, func(t time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}
