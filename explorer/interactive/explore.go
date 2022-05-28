package interactive

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	// "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/jdginn/durins-door/explorer"
	"github.com/jdginn/dwarf-explore/explorer/interactive/style"
)

func getEntryNames(e *explorer.Explorer) ([]list.Item, error) {
	names, err := e.ShowAllChildren()
	items := make([]list.Item, 0, len(names))
	if err != nil {
		return items, err
	}
	for _, n := range names {
		items = append(items, style.ListItem(n))
	}
	return items, nil
}

func stringsToItems(s []string) []list.Item {
	items := make([]list.Item, 0, len(s))
	for _, n := range s {
		items = append(items, style.ListItem(n))
	}
	return items
}

func initExplore(m *model) {
	items, err := getEntryNames(m.explorer)
	if err != nil {
		panic(err)
	}
	m.list = style.BuildList(items, "Select an entry...")
	m.state = explore
}

func ExploreUpdate(m model, msg tea.Msg) (model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	m.textInput, cmd = m.textInput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := m.list.SelectedItem().(style.ListItem)
			if ok {
				m.explorer.StepIntoChild(string(i))
				items, err := getEntryNames(m.explorer)
				if err != nil {
					panic(err)
				}
				m.list = style.BuildList(items, fmt.Sprintf("Currently viewing %s", m.explorer.CurrEntryName()))
			}
			return m, cmd

		// Show info about this entry
		case "i":
			m.msg = m.explorer.Info()

		// Show this entry's type
		case "t":
			proxy, err := m.explorer.GetTypeDefProxy()
			if err != nil {
				panic(err)
			}
			names := proxy.ListChildren()
			items := stringsToItems(names)
			m.list = style.BuildList(items, fmt.Sprintf("Currently viewing type %s", m.explorer.CurrEntryName()))
			m.msg = m.explorer.Info()
			return m, cmd

		// Look at variables instead of types
		case "v":

			// Go up one level
		case "u":
			ok := m.explorer.Up()
			if ok {
				items, err := getEntryNames(m.explorer)
				if err != nil {
					panic(err)
				}
				m.list = style.BuildList(items, fmt.Sprintf("Currently viewing %s", m.explorer.CurrEntryName()))
			}
			return m, cmd

		case "ctrl+c":
			m.state = actionList
			return m, tea.Quit

		case "esc":
			m.state = actionList
			return m, cmd
		}
	}

	return m, cmd
}

func ExploreView(m model) string {
	return m.list.View()
}
