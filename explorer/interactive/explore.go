package interactive

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/jdginn/dwarf-explore/explorer/interactive/style"
)

// Supported functionality:
//
// TypeDefs:
//  list: step into children
//  i: show info
//  l: list all instances,
//      then select a variable from the instances
//  u: go to parent
//
// Variables:
//  list: step into children according to typedef
//  i: show info
//  t: get typedef
//  v: get value
//  u: go to parent according to typedef

type exploreMode int

const (
	variable exploreMode = iota
	typedef
)

func stringsToItems(s []string) []list.Item {
	items := make([]list.Item, 0, len(s))
	for _, n := range s {
		items = append(items, style.ListItem(n))
	}
	return items
}

func initExplore(m *model) {
	items := stringsToItems(m.explorer.ListChildren())
	m.list = style.BuildList(items, "Select an item...")
	m.state = explore
}

// Handle keystrokes shared by all explore actions:
//
//  i:      view verbose info on this entry's attributes
//  b:      go back to the last thing we were doing
//  ctrl+c: exit the program
//  esc:    return to the main menu
func sharedUpdate(m model, keypress string) (model, tea.Cmd) {
	var cmd tea.Cmd
	switch keypress {
	case "i":
		m.msg = m.explorer.Info()

	case "b":
		m.explorer.Back()

	case "ctrl+c":
		m.state = actionList
		cmd = tea.Quit

	case "esc":
		m.state = actionList
	}
	return m, cmd
}

func ExploreUpdate(m model, msg tea.Msg) (model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	m.textInput, cmd = m.textInput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		keypress := msg.String()
		switch keypress {
		case "enter":
			i, ok := m.list.SelectedItem().(style.ListItem)
			if ok {
				m.explorer.StepIntoChild(string(i))
				items := stringsToItems(m.explorer.ListChildren())
				m.list = style.BuildList(items, fmt.Sprintf("Currently viewing %s", m.explorer.CurrName()))
			}
			return m, cmd

		// Show this entry's type
		case "t":
			err := m.explorer.GetType()
			if err != nil {
				panic(err)
			}
			items := stringsToItems(m.explorer.ListChildren())
			m.list = style.BuildList(items, fmt.Sprintf("Currently viewing type %s", m.explorer.CurrName()))
			m.msg = m.explorer.Info()
			return m, cmd

			// Go up one level
		case "u":
			ok := m.explorer.Up()
			if ok {
				items := stringsToItems(m.explorer.ListChildren())
				m.list = style.BuildList(items, fmt.Sprintf("Currently viewing %s", m.explorer.CurrName()))
			}
			return m, cmd
		}
		return sharedUpdate(m, keypress)
	}

	return m, cmd
}

func ExploreView(m model) string {
	return m.list.View()
}
