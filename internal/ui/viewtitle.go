package ui

import (
	"fmt"

	"github.com/rivo/tview"
)

// ViewTitleView displays the current view name in a highlighted, centralized format
type ViewTitleView struct {
	*tview.TextView
	theme *Theme
}

// NewViewTitleView creates a new view title view
func NewViewTitleView() *ViewTitleView {
	theme := DefaultTheme()

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(false).
		SetTextAlign(tview.AlignCenter)

	vtv := &ViewTitleView{
		TextView: textView,
		theme:    theme,
	}

	// No border, just text
	textView.SetBorder(false)

	return vtv
}

// SetViewName sets the view name with formatted label
func (vtv *ViewTitleView) SetViewName(viewName string) {
	// Format with bold colored label: [lightblue::b]View:[white] <view name>
	formattedText := fmt.Sprintf("[lightblue::b]View:[white] %s", viewName)
	vtv.SetText(formattedText)
}
