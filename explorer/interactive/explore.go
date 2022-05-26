package interactive

import (
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

	// Always allow us to quit
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
				m.list = style.BuildList(items, "Select an entry...")
			}
			return m, cmd

		case "u":
			m.explorer.Up()
			items, err := getEntryNames(m.explorer)
			if err != nil {
				panic(err)
			}
			m.list = style.BuildList(items, "Select an entry...")

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
