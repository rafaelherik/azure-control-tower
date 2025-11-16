package ui

import (
	"github.com/gdamore/tcell/v2"
)

// Theme defines the color scheme for the application
type Theme struct {
	Primary    tcell.Color
	Secondary  tcell.Color
	Error      tcell.Color
	Success    tcell.Color
	Warning    tcell.Color
	Info       tcell.Color
	Label      tcell.Color
	Text       tcell.Color
	Border     tcell.Color
	Background tcell.Color
}

// DefaultTheme returns the default theme configuration
func DefaultTheme() *Theme {
	return &Theme{
		Primary:    tcell.ColorBlue,
		Secondary:  tcell.ColorGray,
		Error:      tcell.ColorRed,
		Success:    tcell.ColorGreen,
		Warning:    tcell.ColorYellow,
		Info:       tcell.ColorAqua,
		Label:      tcell.ColorLightBlue, // Light blue for labels like "Tenant:", "Subscription:"
		Text:       tcell.ColorWhite,
		Border:     tcell.ColorBlue,
		Background: tcell.ColorDefault,
	}
}

// GetLabelStyle returns a style for labels (bold, light blue)
func (t *Theme) GetLabelStyle() tcell.Style {
	return tcell.StyleDefault.
		Foreground(t.Label).
		Attributes(tcell.AttrBold)
}

// GetTextStyle returns a style for regular text
func (t *Theme) GetTextStyle() tcell.Style {
	return tcell.StyleDefault.
		Foreground(t.Text)
}

// GetBorderStyle returns a style for borders
func (t *Theme) GetBorderStyle() tcell.Style {
	return tcell.StyleDefault.
		Foreground(t.Border)
}

// GetErrorStyle returns a style for error messages
func (t *Theme) GetErrorStyle() tcell.Style {
	return tcell.StyleDefault.
		Foreground(t.Error)
}

// GetSuccessStyle returns a style for success messages
func (t *Theme) GetSuccessStyle() tcell.Style {
	return tcell.StyleDefault.
		Foreground(t.Success)
}

// GetWarningStyle returns a style for warning messages
func (t *Theme) GetWarningStyle() tcell.Style {
	return tcell.StyleDefault.
		Foreground(t.Warning)
}

// GetInfoStyle returns a style for info messages
func (t *Theme) GetInfoStyle() tcell.Style {
	return tcell.StyleDefault.
		Foreground(t.Info)
}
