package main

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jdginn/durins-door/explorer"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

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

func initialModel() model {
	actions := []list.Item{
		item("Info"),
		item("Get Reader"),
		item("Set Client"),
		item("Get Object"),
		item("List CompileUnits"),
	}

	const defaultWidth = 20

	l := list.New(actions, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Select an action:"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		explorer.NewExplorer(),
		actionList,
		l,
		ti,
		true,
		nil,
	}
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
				i, ok := m.list.SelectedItem().(item)
				if ok {
					m.state = stateMap[string(i)]
					fmt.Printf("State: %s\n", string(i))
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
  if err != nil { panic(err) }
  for _, cu := range CUs {
      s += "\t" + cu + "\n"
    }
	}

	return s
}

func main() {

	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
