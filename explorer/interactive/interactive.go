package interactive

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/jdginn/durins-door/explorer"
	"github.com/jdginn/dwarf-explore/explorer/interactive/style"
)

// Represents the state of the CLI
type state int

const (
	actionList state = iota
	info
	explore
	listCUs
	viewObj
)

func (s state) String() string {
	switch s {
	case actionList:
		return "actionList"
	case info:
		return "info"
	case explore:
		return "explore"
	case listCUs:
		return "listCUs"
	}
	return "unknown"
}

func stateMap() map[string]state {
	return map[string]state{
		"ActionList": actionList,
		"Info":       info,
		"Explore!":   explore,
	}
}

func stateCallbacks() map[string]func(m *model) {
	return map[string]func(m *model){
		"Info":     func(m *model) { m.state = info },
		"Explore!": initExplore,
	}
}

type model struct {
	explorer   *explorer.Explorer
	state      state
	list       list.Model
	textInput  textinput.Model
	textPrompt bool
	err        error
}

func initialModel(file string) model {
	actions := []list.Item{
		style.ListItem("Info"),
		style.ListItem("Explore!"),
	}

	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	m := model{
		explorer.NewExplorer(),
		actionList,
		style.BuildList(actions, "Select an action:"),
		ti,
		true,
		nil,
	}
	m.explorer.CreateReaderFromFile(file)
	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {

	// Basic case to allow selcting an action
	case actionList:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				i, ok := m.list.SelectedItem().(style.ListItem)
				if ok {
					stateCallbacks()[string(i)](&m)
				}
				return m, cmd
			}
		}
		m.list, cmd = m.list.Update(msg)

		m.textInput, cmd = m.textInput.Update(msg)

		// Actions are implemented here
	case listCUs:
		m.state = actionList

	case explore:
		m, msg = ExploreUpdate(m, msg)
	}

	// Always allow us to quit
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			m.state = actionList
		}
	}
	return m, cmd
}

func (m model) View() string {
	s := ""

	switch m.state {
	case actionList:
		s += m.list.View()
	case info:
		s += fmt.Sprintf("Dwarf Explorer:\n")
		s += fmt.Sprintf("\tReader: %s\n", m.explorer.GetReaderFilename())
	case explore:
		s += ExploreView(m)
	case listCUs:
		s += "CUs:\n"
		CUs, err := m.explorer.ListCUs()
		if err != nil {
			panic(err)
		}
		for _, cu := range CUs {
			s += "\t" + cu + "\n"
		}
	}

	return s
}

func Start(file string) {
	p := tea.NewProgram(initialModel(file))
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
