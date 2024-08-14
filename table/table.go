package table

import (
	"encoding/json"
	"fmt"
	prettytable "github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"io"
	"os"
	"slices"
)

type Table struct {
	tableWriter prettytable.Writer
	columns     []string
}

func NewTable(columns []any) (t *Table) {
	t = &Table{
		tableWriter: prettytable.NewWriter(),
	}

	columnsAsStrings := make([]string, 0, len(columns))

	for _, column := range columns {
		columnsAsStrings = append(columnsAsStrings, fmt.Sprintf("%v", column))
	}

	// Sort the columns
	t.columns = columnsAsStrings
	slices.Sort(t.columns)

	// Convert back to any
	var columnsAsAny []any

	for _, column := range t.columns {
		columnsAsAny = append(columnsAsAny, column)
	}

	t.tableWriter.SetStyle(prettytable.StyleColoredBlueWhiteOnBlack)
	t.tableWriter.AppendHeader(columnsAsAny)
	t.tableWriter.SetRowPainter(func(row prettytable.Row) text.Colors {
		return text.Colors{text.FgHiWhite}
	})

	return
}

func NewTableFromJSONTemplate(data json.RawMessage) (t *Table, err error) {
	var mappedData map[string]any
	if err = json.Unmarshal(data, &mappedData); err != nil {
		return
	}

	columns := make([]any, 0, len(mappedData))

	for k := range mappedData {
		columns = append(columns, k)
	}

	t = NewTable(columns)
	return
}

func (t *Table) AddRow(row []any) {
	t.tableWriter.AppendRow(row)
}

func (t *Table) AddRowFromJSON(data json.RawMessage) error {
	var row map[string]any
	err := json.Unmarshal(data, &row)
	if err != nil {
		return err
	}

	var rowSlice []any
	for _, column := range t.columns {
		rowSlice = append(rowSlice, row[fmt.Sprintf("%v", column)])
	}

	t.AddRow(rowSlice)
	return nil
}

func (t Table) Print() {
	t.PrintToWriter(os.Stdout)
}

func (t Table) PrintToWriter(finalW io.Writer) {
	t.tableWriter.SetOutputMirror(finalW)
	t.tableWriter.Render()
}
