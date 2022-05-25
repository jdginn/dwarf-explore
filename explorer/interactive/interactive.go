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

type errMsg error

type model struct {
	explorer   *explorer.Explorer
	state      state
	list       list.Model
	textInput  textinput.Model
	textPrompt bool
	err        error
}

type state int

const (
	actionList state = iota
	info
	getReader
	setClient
	getObj
	listCUs
	viewObj
)

func (s state) String() string {
	switch s {
	case actionList:
		return "actionList"
	case info:
		return "info"
	case getReader:
		return "getReader"
	case setClient:
		return "setClient"
	case getObj:
		return "getObj"
	case listCUs:
		return "listCUs"
	case viewObj:
		return "viewObj"
	}
	return "unknown"
}

func initialModel(file string) model {
	actions := []list.Item{
		style.ListItem("Info"),
		style.ListItem("Explore!"),
		style.ListItem("List CompileUnits"),
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
	case actionList:
		stateMap := map[string]state{
			"ActionList":        actionList,
			"Info":              info,
			"Get Reader":        getReader,
			"Set Client":        setClient,
			"Get Object":        getObj,
			"List CompileUnits": listCUs,
			"View Object":       viewObj,
		}
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				i, ok := m.list.SelectedItem().(style.ListItem)
				if ok {
					m.state = stateMap[string(i)]
				}
				return m, cmd
			}
		}
		m.list, cmd = m.list.Update(msg)

	case getReader:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				m.explorer.CreateReaderFromFile(m.textInput.Value())
				m.state = actionList
			}
		}
		m.textInput, cmd = m.textInput.Update(msg)

	case listCUs:
		m.state = actionList
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
	case getReader:
		s += "Enter the path to a Dwarf Debug file.\n"
		s += m.textInput.View()
	case getObj:
		s += "Enter the path to an object to read from the DWARF.\n"
		s += m.textInput.View()
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
