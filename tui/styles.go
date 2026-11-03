package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// UI Color Palette (Electric Neon / Cyberpunk vibe)
var (
	ColorPrimary   = lipgloss.Color("#00F0FF") // Electric Cyan
	ColorSecondary = lipgloss.Color("#BD00FF") // Neon Violet
	ColorAccent    = lipgloss.Color("#FF007F") // Hot Pink
	ColorSuccess   = lipgloss.Color("#00FF87") // Emerald Neon
	ColorWarning   = lipgloss.Color("#FFE600") // Neon Amber
	ColorMuted     = lipgloss.Color("#6272A4") // Muted Purple-Gray
	ColorSubtle    = lipgloss.Color("#44475A") // Dark Border
	ColorHighlight = lipgloss.Color("#282A36") // Dark Card BG
	ColorWhite     = lipgloss.Color("#F8F8F2") // Off-white
)

// UI Component Styles
var (
	// Header
	HeaderTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#BD00FF")).
				Padding(0, 2).
				MarginRight(1)

	HeaderSubtitleStyle = lipgloss.NewStyle().
				Foreground(ColorPrimary).
				Bold(true)

	HeaderTaglineStyle = lipgloss.NewStyle().
				Foreground(ColorMuted).
				Italic(true)

	HeaderContainer = lipgloss.NewStyle().
			MarginBottom(1).
			PaddingBottom(1).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(ColorSubtle)

	// List Styles
	ListContainerStyle = lipgloss.NewStyle().
				Padding(0, 1)

	CursorStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(ColorPrimary).
				Bold(true)

	UnselectedItemStyle = lipgloss.NewStyle().
				Foreground(ColorWhite)

	MutedTextStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	WhiteTextStyle = lipgloss.NewStyle().
			Foreground(ColorWhite)

	FocusedItemStyle = lipgloss.NewStyle().
				Background(ColorHighlight).
				PaddingLeft(1).
				PaddingRight(1)

	CheckboxSelectedStyle = lipgloss.NewStyle().
				Foreground(ColorSuccess).
				Bold(true)

	CheckboxUnselectedStyle = lipgloss.NewStyle().
				Foreground(ColorMuted)

	// Preview / Details Card
	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSecondary).
			Padding(1, 2).
			MarginLeft(1)

	CardTitleStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true).
			MarginBottom(1)

	BadgeKeyStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	BadgeValStyle = lipgloss.NewStyle().
			Foreground(ColorWhite)

	TagStyle = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Background(ColorSecondary).
			Padding(0, 1).
			MarginRight(1).
			Bold(true)

	// Dialog / Modal
	ModalStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorWarning).
			Padding(1, 3).
			MarginTop(1).
			MarginBottom(1)

	ButtonActive = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(ColorPrimary).
			Bold(true).
			Padding(0, 2).
			MarginRight(2)

	ButtonInactive = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Background(ColorSubtle).
			Padding(0, 2).
			MarginRight(2)

	// Footer / Shortcuts
	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			MarginRight(2)

	FooterContainer = lipgloss.NewStyle().
			MarginTop(1).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(ColorSubtle).
			PaddingTop(1)

	// Success / Summary View
	SuccessBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSuccess).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	CommandBox = lipgloss.NewStyle().
			Background(ColorHighlight).
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorPrimary).
			Padding(0, 1).
			MarginTop(1)
)
