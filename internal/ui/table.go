package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ColumnConfig defines a table column configuration
type ColumnConfig struct {
	Name       string
	Width      int // 0 means auto-width
	Align      int // tview.AlignLeft, tview.AlignCenter, tview.AlignRight
	Selectable bool
}

// RowAction represents an action that can be performed on a row
type RowAction struct {
	Key      tcell.Key
	Rune     rune
	Label    string
	Callback func(rowIndex int, data interface{}) bool // Returns true if event was handled
}

// TableConfig holds the configuration for a table view
type TableConfig struct {
	Title        string
	Columns      []ColumnConfig
	RowActions   []RowAction
	OnSelect     func(rowIndex int, data interface{})
	GetRowData   func(rowIndex int) interface{}                 // Function to get row data by index
	GetCellValue func(data interface{}, columnIndex int) string // Function to extract cell value from row data
}

// TableView is a reusable table component that accepts header configuration and row actions
type TableView struct {
	*tview.Table
	config          *TableConfig
	data            []interface{} // Store row data
	filterText      string        // Current filter text
	filteredIndices []int         // Indices of filtered rows
	theme           *Theme
}

// NewTableView creates a new table view with the given configuration
func NewTableView(config *TableConfig) *TableView {
	theme := DefaultTheme()

	table := tview.NewTable().
		SetSelectable(true, false).
		SetFixed(1, 0)

	tv := &TableView{
		Table:           table,
		config:          config,
		filteredIndices: []int{},
		theme:           theme,
	}

	// Always enable border with theme colors
	table.SetBorder(true).
		SetBorderColor(theme.Border)

	// Set title if provided (this will be the top border with actions)
	if config.Title != "" {
		table.SetTitle(config.Title)
	}

	// Set up selection callback
	table.SetSelectedFunc(func(row, column int) {
		if row > 0 && config.OnSelect != nil {
			dataIndex := tv.getDataIndex(row - 1)
			if dataIndex >= 0 && dataIndex < len(tv.data) {
				config.OnSelect(dataIndex, tv.data[dataIndex])
			}
		}
	})

	// Render headers
	tv.renderHeaders()

	return tv
}

// SetConfig updates the table configuration
func (tv *TableView) SetConfig(config *TableConfig) {
	tv.config = config
	// Ensure border is always enabled
	tv.SetBorder(true).
		SetBorderColor(tv.theme.Border)

	// Set title if provided
	if config.Title != "" {
		tv.SetTitle(config.Title)
	}
	tv.renderHeaders()
	tv.RenderData()
	// expandColumns is already called in RenderData
}

// LoadData loads data into the table
func (tv *TableView) LoadData(data []interface{}) {
	tv.data = data
	tv.filteredIndices = make([]int, len(data))
	for i := range data {
		tv.filteredIndices[i] = i
	}
	tv.RenderData()
}

// RenderData renders all data rows
func (tv *TableView) RenderData() {
	// Clear existing data rows (keep header)
	rowCount := tv.GetRowCount()
	for i := rowCount - 1; i > 0; i-- {
		tv.RemoveRow(i)
	}

	// Render filtered rows
	for i, dataIndex := range tv.filteredIndices {
		if dataIndex < 0 || dataIndex >= len(tv.data) {
			continue
		}
		data := tv.data[dataIndex]
		for colIndex, col := range tv.config.Columns {
			cellValue := tv.config.GetCellValue(data, colIndex)
			cell := tview.NewTableCell(cellValue).
				SetTextColor(tv.theme.Text).
				SetExpansion(1) // Make cells expand to fill available space
			if col.Align != 0 {
				cell.SetAlign(col.Align)
			}
			tv.SetCell(i+1, colIndex, cell)
		}
	}

	// Expand columns to use full width
	tv.expandColumns()
}

// renderHeaders renders the header row
func (tv *TableView) renderHeaders() {
	for i, col := range tv.config.Columns {
		cell := tview.NewTableCell(col.Name).
			SetSelectable(false).
			SetAttributes(tcell.AttrBold).
			SetTextColor(tv.theme.Label).
			SetExpansion(1) // Make header cells expand to fill available space
		if col.Align != 0 {
			cell.SetAlign(col.Align)
		}
		tv.SetCell(0, i, cell)
	}
}

// expandColumns expands all columns to use the full available width
func (tv *TableView) expandColumns() {
	numColumns := len(tv.config.Columns)
	if numColumns == 0 {
		return
	}

	// Set all columns to expand equally
	// In tview, columns with expansion > 0 will share available space
	// We set expansion to 1 for all columns so they share equally
	for colIndex := 0; colIndex < numColumns; colIndex++ {
		// Update header cell expansion
		if cell := tv.GetCell(0, colIndex); cell != nil {
			cell.SetExpansion(1)
		}

		// Update all data row cells for this column
		for row := 1; row < tv.GetRowCount(); row++ {
			if cell := tv.GetCell(row, colIndex); cell != nil {
				cell.SetExpansion(1)
			}
		}
	}
}

// SetFilter sets the filter text and updates the view
func (tv *TableView) SetFilter(filterText string) {
	tv.filterText = filterText
	if filterText == "" {
		// Show all rows
		tv.filteredIndices = make([]int, len(tv.data))
		for i := range tv.data {
			tv.filteredIndices[i] = i
		}
	} else {
		// Filter rows based on cell values
		tv.filteredIndices = []int{}
		for i, data := range tv.data {
			matches := false
			for colIndex := range tv.config.Columns {
				cellValue := tv.config.GetCellValue(data, colIndex)
				if containsIgnoreCase(cellValue, filterText) {
					matches = true
					break
				}
			}
			if matches {
				tv.filteredIndices = append(tv.filteredIndices, i)
			}
		}
	}
	tv.RenderData()
	// Note: Selection will be maintained by tview automatically
}

// GetFilter returns the current filter text
func (tv *TableView) GetFilter() string {
	return tv.filterText
}

// ClearFilter clears the filter
func (tv *TableView) ClearFilter() {
	tv.SetFilter("")
}

// HandleKey handles key events for the table
func (tv *TableView) HandleKey(event *tcell.EventKey) *tcell.EventKey {
	row, _ := tv.GetSelection()
	if row > 0 {
		dataIndex := tv.getDataIndex(row - 1)
		if dataIndex >= 0 && dataIndex < len(tv.data) {
			// Check row actions
			for _, action := range tv.config.RowActions {
				if action.Key != 0 && event.Key() == action.Key {
					if action.Callback != nil && action.Callback(dataIndex, tv.data[dataIndex]) {
						return nil
					}
				}
				if action.Rune != 0 && event.Key() == tcell.KeyRune && event.Rune() == action.Rune {
					if action.Callback != nil && action.Callback(dataIndex, tv.data[dataIndex]) {
						return nil
					}
				}
			}
		}
	}
	return event
}

// getDataIndex converts a display row index to a data index
func (tv *TableView) getDataIndex(displayRowIndex int) int {
	if displayRowIndex < 0 || displayRowIndex >= len(tv.filteredIndices) {
		return -1
	}
	return tv.filteredIndices[displayRowIndex]
}

// GetSelectedData returns the data for the currently selected row
func (tv *TableView) GetSelectedData() interface{} {
	row, _ := tv.GetSelection()
	if row > 0 {
		dataIndex := tv.getDataIndex(row - 1)
		if dataIndex >= 0 && dataIndex < len(tv.data) {
			return tv.data[dataIndex]
		}
	}
	return nil
}

// GetRowCount returns the number of data rows (excluding header)
func (tv *TableView) GetDataRowCount() int {
	return len(tv.filteredIndices)
}

// SetTopTitle sets the top border title with actions formatted as buttons
func (tv *TableView) SetTopTitle(actions string) {
	// Format actions as buttons
	formattedActions := formatActionsAsButtons(actions)
	tv.SetTitle(formattedActions)
}

// containsIgnoreCase checks if a string contains a substring (case-insensitive)
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) &&
		(len(substr) == 0 ||
			containsIgnoreCaseHelper(s, substr))
}

func containsIgnoreCaseHelper(s, substr string) bool {
	sLower := toLower(s)
	substrLower := toLower(substr)
	for i := 0; i <= len(sLower)-len(substrLower); i++ {
		if sLower[i:i+len(substrLower)] == substrLower {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	result := make([]rune, len(s))
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			result[i] = r + 32
		} else {
			result[i] = r
		}
	}
	return string(result)
}
