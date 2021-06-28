package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Table struct {
	Header []string
	Rows   [][]interface{}
}

func NewTable(header []string) *Table {
	return &Table{
		Header: header,
		Rows:   make([][]interface{}, 0),
	}
}

func (t *Table) AddRow(row []interface{}) {
	t.Rows = append(t.Rows, row)
}

func (t *Table) Format(w io.Writer) error {
	var buf bytes.Buffer

	rows := t.Render()
	widths := t.ColumnWidths(rows)

	for i, label := range t.Header {
		if i > 0 {
			buf.WriteString("  ")
		}

		fmt.Fprintf(&buf, "%-*s", widths[i], strings.ToUpper(label))
	}
	buf.WriteByte('\n')

	for _, row := range rows {
		for j, s := range row {
			if j > 0 {
				buf.WriteString("  ")
			}

			fmt.Fprintf(&buf, "%-*s", widths[j], s)
		}

		buf.WriteByte('\n')
	}

	_, err := io.Copy(w, &buf)
	return err
}

func (t *Table) Render() [][]string {
	rows := make([][]string, len(t.Rows))

	for i, row := range t.Rows {
		rows[i] = make([]string, len(row))

		for j, value := range row {
			rows[i][j] = t.RenderValue(value)
		}
	}

	return rows
}

func (t *Table) RenderValue(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

func (t *Table) ColumnWidths(rows [][]string) []int {
	widths := make([]int, len(t.Header))

	for i, label := range t.Header {
		widths[i] = len(label)
	}

	for _, row := range rows {
		for j, value := range row {
			if len(value) > widths[j] {
				widths[j] = len(value)
			}
		}
	}

	return widths
}
