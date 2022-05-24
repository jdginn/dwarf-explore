package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/jdginn/durins-door/explorer"
)

type errMsg error

type state int

const (
	getReader state = iota
	getClient
	getObjPath
)

type model struct {
	explorer   *explorer.Explorer
	state      state
	textInput  textinput.Model
	textPrompt bool
	err        error
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		explorer.NewExplorer(),
		getReader,
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

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			switch m.state {

			case getReader:
				m.explorer.GetReaderFromFile(m.textInput.Value())
				m.state = getObjPath
			}
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	switch m.state {
	case getReader:
		m.textInput, cmd = m.textInput.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	s := ""

	switch m.state {
	case getReader:
		s += "Welcome to dwarf-explore\n"
		s += "Enter the path to a Dwarf Debug file.\n"
		s += "Press q to quit.\n"
  case getObjPath:
    s += "Enter the path to an object to read from the DWARF.\n"
		s += "Press q to quit.\n"
	default:
		s += "Press q to quit. \n"
	}

	if m.textPrompt {
		s += m.textInput.View()
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
