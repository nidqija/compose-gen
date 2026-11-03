package tui

import (
	"embed"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type appState int

const (
	stateSelect appState = iota
	stateConfirmOverwrite
	stateDone
)

type confirmChoice int

const (
	choiceMerge confirmChoice = iota
	choiceOverwrite
	choiceCancel
)

// Model represents the Bubble Tea UI state
type Model struct {
	fs             embed.FS
	services       []ServiceItem
	cursor         int
	state          appState
	confirmChoice  confirmChoice
	width          int
	height         int
	log            GenerationLog
	composeExisted bool
	errorMessage   string
}

// NewModel initializes the TUI
func NewModel(fs embed.FS) (Model, error) {
	services, err := LoadServices(fs)
	if err != nil {
		return Model{}, err
	}

	_, composeErr := os.Stat("docker-compose.yml")
	composeExisted := composeErr == nil

	return Model{
		fs:             fs,
		services:       services,
		cursor:         0,
		state:          stateSelect,
		confirmChoice:  choiceMerge,
		width:          80,
		height:         24,
		composeExisted: composeExisted,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q", "esc":
			if m.state == stateConfirmOverwrite {
				m.state = stateSelect
				return m, nil
			}
			return m, tea.Quit
		}

		switch m.state {
		case stateSelect:
			return m.updateSelect(msg)
		case stateConfirmOverwrite:
			return m.updateConfirm(msg)
		case stateDone:
			if msg.String() == "enter" || msg.String() == "q" || msg.String() == "esc" {
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m Model) updateSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.errorMessage = ""

	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		} else {
			m.cursor = len(m.services) - 1
		}

	case "down", "j":
		if m.cursor < len(m.services)-1 {
			m.cursor++
		} else {
			m.cursor = 0
		}

	case " ", "x":
		if len(m.services) > 0 {
			m.services[m.cursor].Selected = !m.services[m.cursor].Selected
		}

	case "a":
		// Check if all are selected
		allSelected := true
		for _, s := range m.services {
			if !s.Selected {
				allSelected = false
				break
			}
		}
		// Toggle all
		for i := range m.services {
			m.services[i].Selected = !allSelected
		}

	case "enter":
		selected := m.getSelectedServices()
		if len(selected) == 0 {
			m.errorMessage = "Please select at least one service with [Space] before generating."
			return m, nil
		}

		if m.composeExisted {
			m.state = stateConfirmOverwrite
			return m, nil
		}

		// Generate immediately
		m.log = GenerateCompose(m.fs, selected, false)
		m.state = stateDone
		return m, nil
	}

	return m, nil
}

func (m Model) updateConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "left", "h", "up":
		if m.confirmChoice > 0 {
			m.confirmChoice--
		} else {
			m.confirmChoice = choiceCancel
		}

	case "right", "l", "down", "tab":
		if m.confirmChoice < choiceCancel {
			m.confirmChoice++
		} else {
			m.confirmChoice = choiceMerge
		}

	case "enter":
		switch m.confirmChoice {
		case choiceMerge:
			selected := m.getSelectedServices()
			m.log = GenerateCompose(m.fs, selected, true)
			m.state = stateDone

		case choiceOverwrite:
			selected := m.getSelectedServices()
			m.log = GenerateCompose(m.fs, selected, false)
			m.state = stateDone

		case choiceCancel:
			m.state = stateSelect
		}
	}

	return m, nil
}

func (m Model) getSelectedServices() []ServiceItem {
	var sel []ServiceItem
	for _, s := range m.services {
		if s.Selected {
			sel = append(sel, s)
		}
	}
	return sel
}

func (m Model) View() string {
	switch m.state {
	case stateConfirmOverwrite:
		return m.viewConfirm()
	case stateDone:
		return m.viewDone()
	default:
		return m.viewSelect()
	}
}

func renderBigTitle(width int) string {
	if width >= 76 {
		lines := []string{
			`   ______                                            ______           `,
			`  / ____/___  ____ ___  ____  ____  ________        / ____/__  ____   `,
			` / /   / __ \/ __ ' _ \/ __ \/ __ \/ ___/ _ \______/ / __/ _ \/ __ \  `,
			`/ /___/ /_/ / / / / / / /_/ / /_/ (__  )  __/_____/ /_/ /  __/ / / /  `,
			`\____/\____/_/ /_/ /_/ .___/\____/____/\___/      \____/\___/_/ /_/   `,
			`                    /_/                                               `,
		}
		colors := []string{
			"#00F0FF",
			"#30D5FF",
			"#6DB2FF",
			"#998FFF",
			"#BD00FF",
			"#D900FF",
		}
		var styledLines []string
		for i, line := range lines {
			styledLines = append(styledLines, lipgloss.NewStyle().Foreground(lipgloss.Color(colors[i])).Bold(true).Render(line))
		}
		return strings.Join(styledLines, "\n")
	}

	// Compact version for narrower terminals (52 cols)
	lines := []string{
		`  ___ ___  __  __ ___  ___  ___ ___     ___ ___ _  _ `,
		` / __/ _ \|  \/  | _ \/ _ \/ __| __|   / __| __| \| |`,
		`| (_| (_) | |\/| |  _/ (_) \__ \ _|   | (_ | _|| .  |`,
		` \___\___/|_|  |_|_|  \___/|___/___|   \___|___|_|\_|`,
	}
	colors := []string{
		"#00F0FF",
		"#58CAFF",
		"#998FFF",
		"#BD00FF",
	}
	var styledLines []string
	for i, line := range lines {
		styledLines = append(styledLines, lipgloss.NewStyle().Foreground(lipgloss.Color(colors[i])).Bold(true).Render(line))
	}
	return strings.Join(styledLines, "\n")
}

func (m Model) viewSelect() string {
	var b strings.Builder

	// Big Title Header
	bigTitle := renderBigTitle(m.width)
	subtitle := HeaderSubtitleStyle.Render("⚡ Modern Docker Stack Generator ⚡")
	tagline := HeaderTaglineStyle.Render("Curate your container ecosystem with ease")

	headerBlock := lipgloss.JoinVertical(lipgloss.Left,
		bigTitle,
		"",
		subtitle,
		tagline,
	)
	b.WriteString(HeaderContainer.Render(headerBlock))
	b.WriteString("\n")

	// Count badge
	selectedCount := 0
	for _, s := range m.services {
		if s.Selected {
			selectedCount++
		}
	}

	statusLine := fmt.Sprintf(" Available Services (%d/%d selected):", selectedCount, len(m.services))
	b.WriteString(BadgeKeyStyle.Render(statusLine))
	b.WriteString("\n\n")

	// Render Service List (Left column)
	var listItems []string
	for i, item := range m.services {
		cursor := "  "
		isCurrent := i == m.cursor
		if isCurrent {
			cursor = CursorStyle.Render("❯ ")
		}

		checkbox := CheckboxUnselectedStyle.Render("[ ]")
		if item.Selected {
			checkbox = CheckboxSelectedStyle.Render("[✓]")
		}

		label := UnselectedItemStyle.Render(item.Label)
		if item.Selected {
			label = SelectedItemStyle.Render(item.Label)
		}

		line := fmt.Sprintf("%s%s %s", cursor, checkbox, label)
		if isCurrent {
			line = FocusedItemStyle.Render(line)
		}
		listItems = append(listItems, line)
	}

	leftPane := ListContainerStyle.Render(strings.Join(listItems, "\n"))

	// Render Details Card (Right column)
	var rightPane string
	if len(m.services) > 0 && m.cursor < len(m.services) {
		current := m.services[m.cursor]
		var card strings.Builder

		card.WriteString(CardTitleStyle.Render(fmt.Sprintf("📦 %s", current.Key)))
		card.WriteString("\n\n")

		if current.Image != "" {
			card.WriteString(BadgeKeyStyle.Render("Image: ") + BadgeValStyle.Render(current.Image) + "\n")
		}

		if len(current.Ports) > 0 {
			card.WriteString(BadgeKeyStyle.Render("Ports: ") + BadgeValStyle.Render(strings.Join(current.Ports, ", ")) + "\n")
		}

		if len(current.Volumes) > 0 {
			var volNames []string
			for k := range current.Volumes {
				volNames = append(volNames, k)
			}
			card.WriteString(BadgeKeyStyle.Render("Volumes: ") + BadgeValStyle.Render(strings.Join(volNames, ", ")) + "\n")
		}

		if len(current.ExtraFiles) > 0 {
			card.WriteString("\n" + BadgeKeyStyle.Render("Scaffolds:") + "\n")
			for _, f := range current.ExtraFiles {
				card.WriteString(BadgeValStyle.Render("  • "+f) + "\n")
			}
		}

		if len(current.EnvVars) > 0 {
			card.WriteString("\n" + BadgeKeyStyle.Render("Env Config:") + "\n")
			limit := len(current.EnvVars)
			if limit > 4 {
				limit = 4
			}
			for i := 0; i < limit; i++ {
				card.WriteString(BadgeValStyle.Render("  • "+current.EnvVars[i]) + "\n")
			}
			if len(current.EnvVars) > 4 {
				card.WriteString(MutedTextStyle.Render(fmt.Sprintf("  ... (+%d more)", len(current.EnvVars)-4)) + "\n")
			}
		}

		rightPane = CardStyle.Width(38).Render(card.String())
	}

	// Layout side-by-side or stacked based on terminal width
	var body string
	if m.width > 85 {
		body = lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
	} else {
		body = lipgloss.JoinVertical(lipgloss.Left, leftPane, rightPane)
	}
	b.WriteString(body)

	// Error message if any
	if m.errorMessage != "" {
		b.WriteString("\n\n" + lipgloss.NewStyle().Foreground(ColorAccent).Bold(true).Render("⚠️  "+m.errorMessage))
	}

	// Footer Help
	helpKeys := []string{
		HelpKeyStyle.Render("↑/↓, j/k") + " " + HelpDescStyle.Render("Navigate"),
		HelpKeyStyle.Render("Space") + " " + HelpDescStyle.Render("Toggle"),
		HelpKeyStyle.Render("a") + " " + HelpDescStyle.Render("All"),
		HelpKeyStyle.Render("Enter") + " " + HelpDescStyle.Render("Generate"),
		HelpKeyStyle.Render("q") + " " + HelpDescStyle.Render("Quit"),
	}
	footer := FooterContainer.Render(strings.Join(helpKeys, " "))
	b.WriteString("\n" + footer)

	return b.String()
}

func (m Model) viewConfirm() string {
	var b strings.Builder

	b.WriteString(HeaderTitleStyle.Render("⚠️  EXISTING FILE DETECTED") + "\n\n")

	content := lipgloss.NewStyle().Foreground(ColorWhite).Render(
		"A docker-compose.yml already exists in this directory.\n" +
			"How would you like compose-gen to proceed?\n\n",
	)
	b.WriteString(content)

	btnMerge := ButtonInactive.Render("Merge Services")
	btnOverwrite := ButtonInactive.Render("Overwrite Completely")
	btnCancel := ButtonInactive.Render("Cancel")

	switch m.confirmChoice {
	case choiceMerge:
		btnMerge = ButtonActive.Render("Merge Services")
	case choiceOverwrite:
		btnOverwrite = ButtonActive.Render("Overwrite Completely")
	case choiceCancel:
		btnCancel = ButtonActive.Render("Cancel")
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, btnMerge, btnOverwrite, btnCancel)
	b.WriteString(ModalStyle.Render(buttons))
	b.WriteString("\n\n")

	help := HelpKeyStyle.Render("←/→, Tab") + " " + HelpDescStyle.Render("Switch Option") + "  " +
		HelpKeyStyle.Render("Enter") + " " + HelpDescStyle.Render("Confirm") + "  " +
		HelpKeyStyle.Render("Esc") + " " + HelpDescStyle.Render("Back")
	b.WriteString(FooterContainer.Render(help))

	return b.String()
}

func (m Model) viewDone() string {
	var b strings.Builder

	header := HeaderTitleStyle.Render("✨ GENERATION COMPLETE ✨")
	b.WriteString(header + "\n\n")

	if m.log.Error != nil {
		errBox := lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true).
			Render("❌ Error: " + m.log.Error.Error())
		b.WriteString(errBox + "\n")
		return b.String()
	}

	var summary strings.Builder
	modeText := "Created new docker-compose.yml"
	if m.log.Mode == "merged" {
		modeText = "Merged into existing docker-compose.yml"
	} else if m.log.Mode == "overwritten" {
		modeText = "Overwrote docker-compose.yml"
	}

	summary.WriteString(lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true).Render("✔ "+modeText) + "\n\n")

	summary.WriteString(BadgeKeyStyle.Render("Services Activated:") + "\n")
	for _, s := range m.log.ServicesGenerated {
		summary.WriteString(WhiteTextStyle.Render(fmt.Sprintf("  • %s\n", s)))
	}

	if len(m.log.FilesScaffolded) > 0 {
		summary.WriteString("\n" + BadgeKeyStyle.Render("Scaffolded Files:") + "\n")
		for _, f := range m.log.FilesScaffolded {
			summary.WriteString(lipgloss.NewStyle().Foreground(ColorPrimary).Render(fmt.Sprintf("  ✔ %s\n", f)))
		}
	}

	if len(m.log.ExistingPreserved) > 0 {
		summary.WriteString("\n" + BadgeKeyStyle.Render("Preserved Existing Files:") + "\n")
		for _, f := range m.log.ExistingPreserved {
			summary.WriteString(MutedTextStyle.Render(fmt.Sprintf("  • %s (kept existing)\n", f)))
		}
	}

	b.WriteString(SuccessBox.Render(summary.String()))
	b.WriteString("\n\n")

	// Quickstart commands
	cmdHeader := BadgeKeyStyle.Render("🚀 Next Steps:")
	cmds := fmt.Sprintf(
		"  %s   Start containers in background\n"+
			"  %s   Check running services\n"+
			"  %s   View live logs\n"+
			"  %s   Stop services",
		lipgloss.NewStyle().Foreground(ColorSuccess).Render("compose-gen up"),
		lipgloss.NewStyle().Foreground(ColorSuccess).Render("compose-gen ps"),
		lipgloss.NewStyle().Foreground(ColorSuccess).Render("compose-gen logs"),
		lipgloss.NewStyle().Foreground(ColorSuccess).Render("compose-gen down"),
	)
	b.WriteString(cmdHeader + "\n" + CommandBox.Render(cmds) + "\n\n")

	b.WriteString(FooterContainer.Render(HelpKeyStyle.Render("Press [Enter] or [q] to exit")))

	return b.String()
}
