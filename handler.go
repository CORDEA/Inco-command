package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"strings"
)

const (
	SelectedMark = "[*] "
)

type Handler struct {
	request    *Request
	tableName  string
	dialogName string
	titles     []string
	histories  []History
}

func NewHandler(
	request *Request,
	tableName string,
	dialogName string,
	titles []string,
	histories []History,
) *Handler {
	return &Handler{
		request:    request,
		tableName:  tableName,
		dialogName: dialogName,
		titles:     titles,
		histories:  histories,
	}
}

func (h *Handler) KeyBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding(h.tableName, gocui.KeyArrowDown, gocui.ModNone, h.cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(h.tableName, gocui.KeyArrowUp, gocui.ModNone, h.cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(h.tableName, gocui.KeyEnter, gocui.ModNone, h.selectLine); err != nil {
		return err
	}
	if err := g.SetKeybinding(h.tableName, gocui.KeyCtrlD, gocui.ModNone, h.deleteHistories); err != nil {
		return err
	}

	if err := g.SetKeybinding(h.dialogName, gocui.KeyCtrlY, gocui.ModNone, h.dialogYes); err != nil {
		return err
	}
	if err := g.SetKeybinding(h.dialogName, gocui.KeyCtrlN, gocui.ModNone, h.dialogNo); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, h.quit); err != nil {
		return err
	}
	return nil
}

func (h *Handler) cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}
	x, y := v.Cursor()
	err := v.SetCursor(x, y+1)
	if err == nil {
		return nil
	}
	x, y = v.Origin()
	if err = v.SetOrigin(x, y+1); err != nil {
		return err
	}
	return err
}

func (h *Handler) cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}
	x, y := v.Cursor()
	err := v.SetCursor(x, y-1)
	if err == nil {
		return nil
	}
	x, y = v.Origin()
	if err = v.SetOrigin(x, y-1); err != nil {
		return err
	}
	return err
}

func (h *Handler) selectLine(g *gocui.Gui, v *gocui.View) error {
	_, y := v.Cursor()
	if len(h.titles) <= y {
		return nil
	}

	v.Clear()
	if strings.HasPrefix(h.titles[y], SelectedMark) {
		h.titles[y] = strings.TrimLeft(h.titles[y], SelectedMark)
	} else {
		h.titles[y] = SelectedMark + h.titles[y]
	}

	for i := range h.titles {
		fmt.Fprintln(v, h.titles[i])
	}
	return nil
}

func (h *Handler) deleteHistories(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(h.dialogName, maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Do you want to delete it? (y/n)")
		if _, err := g.SetCurrentView(h.dialogName); err != nil {
			return err
		}
	}
	return nil
}

func (h *Handler) dialogYes(g *gocui.Gui, v *gocui.View) error {
	selected := []History{}
	for i := range h.titles {
		if strings.HasPrefix(h.titles[i], SelectedMark) {
			selected = append(selected, h.histories[i])
		}
	}
	v.Clear()
	if err := h.request.DeleteHistories(selected); err != nil {
		fmt.Fprintln(v, err)
	} else {
		fmt.Fprintln(v, "Completed")
	}
	return nil
}

func (h *Handler) dialogNo(g *gocui.Gui, v *gocui.View) error {
	g.DeleteView(h.dialogName)
	g.SetCurrentView(h.tableName)
	return nil
}

func (h *Handler) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
