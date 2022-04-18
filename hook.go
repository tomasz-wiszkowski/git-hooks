package main

import (
	"fmt"

	"github.com/rivo/tview"
)

type SelectedState int32

const (
	/// Not configured in config file.
	SelectedStateUnknown SelectedState = iota
	/// Not available due to missing dependency.
	SelectedStateUnavailable
	/// Disabled in config file.
	SelectedStateDisabled
	/// Enabled in config file.
	SelectedStateEnabled
)

type Hook struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	selState SelectedState
}

func (h *Hook) toggleSelected() {
	if h.selState == SelectedStateUnknown || h.selState == SelectedStateUnavailable {
		// TODO: check if possible to enable.
		h.selState = SelectedStateEnabled
	} else if h.selState == SelectedStateDisabled {
		h.selState = SelectedStateEnabled
	} else {
		h.selState = SelectedStateDisabled
	}
}

func (h *Hook) updateNode(node *tview.TreeNode) {
	var marker rune
	if h.selState == SelectedStateUnknown {
		marker = '?'
	} else if h.selState == SelectedStateUnavailable {
		marker = '✘'
	} else if h.selState == SelectedStateDisabled {
		marker = ' '
	} else if h.selState == SelectedStateEnabled {
		marker = '✔'
	}
	node.SetText(fmt.Sprintf("[%c] %s", marker, h.Name))
}


