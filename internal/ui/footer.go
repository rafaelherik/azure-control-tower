package ui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

// FooterView displays table statistics and other footer information
type FooterView struct {
	*tview.TextView
	theme *Theme
}

// NewFooterView creates a new footer view
func NewFooterView() *FooterView {
	theme := DefaultTheme()

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(false).
		SetTextAlign(tview.AlignLeft)

	fv := &FooterView{
		TextView: textView,
		theme:    theme,
	}

	// Set border with theme colors
	textView.SetBorder(true).
		SetBorderColor(theme.Border)

	return fv
}

// UpdateCount updates the footer with table item counts
func (fv *FooterView) UpdateCount(totalCount, filteredCount int, hasFilter bool) {
	var text string

	if hasFilter && filteredCount != totalCount {
		// Show both filtered and total count
		text = fmt.Sprintf("[lightblue::b]Items:[white] Showing %d of %d", filteredCount, totalCount)
	} else {
		// Show only total count
		text = fmt.Sprintf("[lightblue::b]Items:[white] %d", totalCount)
	}

	fv.SetText(text)
}

// UpdateCountWithActions updates the footer with counts and action keys
func (fv *FooterView) UpdateCountWithActions(totalCount, filteredCount int, hasFilter bool, actions string) {
	var text string

	if hasFilter && filteredCount != totalCount {
		// Show both filtered and total count
		text = fmt.Sprintf("[lightblue::b]Items:[white] Showing %d of %d  |  %s", filteredCount, totalCount, formatActionsAsButtons(actions))
	} else {
		// Show only total count
		text = fmt.Sprintf("[lightblue::b]Items:[white] %d  |  %s", totalCount, formatActionsAsButtons(actions))
	}

	fv.SetText(text)
}

// formatActionsAsButtons formats action keys as button-like elements
func formatActionsAsButtons(actions string) string {
	// Split actions by comma
	parts := strings.Split(actions, ", ")
	var buttons []string
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		
		// Format as button: [white:blue] action [white]
		// This creates a blue background with white text
		button := fmt.Sprintf("[white:blue] %s [white]", part)
		buttons = append(buttons, button)
	}
	
	// Join with extra spacing between buttons
	return strings.Join(buttons, "  ")
}

// Clear clears the footer
func (fv *FooterView) Clear() {
	fv.SetText("")
}
