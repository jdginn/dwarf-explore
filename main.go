package main

import (
  "fmt"
  "os"

	tea "github.com/charmbracelet/bubbletea"

  "github.com/jdginn/durins-door/explorer"
)

type model struct {
  explorer *explorer.Explorer
}

func initialModel() model {
	return model{nil}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
  s := "Welcome to dwarf-explore\n"
  s += "Enter the path to a Dwarf Debug file.\n"
  s += "Press q to quit.\n"
  return s
}

func main() {
    p := tea.NewProgram(initialModel())
    if err := p.Start(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}
