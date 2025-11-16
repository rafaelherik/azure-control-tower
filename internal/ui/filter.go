package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// FilterMode handles the filter mode (/) input
type FilterMode struct {
	app        *tview.Application
	inputField *tview.InputField
	visible    bool
	onFilter   func(filterText string)
	onCancel   func()
	theme      *Theme
}

// NewFilterMode creates a new filter mode handler
func NewFilterMode(app *tview.Application) *FilterMode {
	theme := DefaultTheme()

	inputField := tview.NewInputField().
		SetLabel("[lightblue::b]/[white]").
		SetFieldWidth(0).
		SetFieldTextColor(theme.Text).
		SetLabelColor(theme.Label)

	fm := &FilterMode{
		app:        app,
		inputField: inputField,
		visible:    false,
		theme:      theme,
	}

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			filterText := inputField.GetText()
			fm.Hide()
			if fm.onFilter != nil {
				fm.onFilter(filterText)
			}
		} else if key == tcell.KeyEsc {
			fm.Hide()
			if fm.onCancel != nil {
				fm.onCancel()
			}
		}
	})

	return fm
}

// Show displays the filter input field
func (fm *FilterMode) Show() {
	fm.visible = true
	fm.inputField.SetText("")
	fm.app.SetFocus(fm.inputField)
}

// Hide hides the filter input field
func (fm *FilterMode) Hide() {
	fm.visible = false
	fm.app.SetFocus(nil)
}

// IsVisible returns whether filter mode is currently visible
func (fm *FilterMode) IsVisible() bool {
	return fm.visible
}

// GetInputField returns the input field for embedding in layouts
func (fm *FilterMode) GetInputField() *tview.InputField {
	return fm.inputField
}

// SetOnFilter sets the callback for when a filter is applied
func (fm *FilterMode) SetOnFilter(callback func(string)) {
	fm.onFilter = callback
}

// SetOnCancel sets the callback for when filter mode is cancelled
func (fm *FilterMode) SetOnCancel(callback func()) {
	fm.onCancel = callback
}
