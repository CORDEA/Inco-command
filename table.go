package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type Table struct {
	tableName string
	items     []string
}

func NewTable(tableName string, items []string) *Table {
	return &Table{
		tableName: tableName,
		items:     items,
	}
}

func (t *Table) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(t.tableName, -1, -1, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for i := range t.items {
			fmt.Fprintln(v, t.items[i])
		}
		if _, err := g.SetCurrentView(t.tableName); err != nil {
			return err
		}
	}
	return nil
}
