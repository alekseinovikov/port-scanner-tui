package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
	"strings"
	"time"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

const (
	padding  = 2
	maxWidth = 80
)

type model struct {
	progress           progress.Model
	textInput          textinput.Model
	openPorts          []uint
	scanned            int
	addressInputActive bool
	err                error
}

type results struct {
	ports    []uint
	finished bool
	scanned  int
}

type (
	errMsg error
)

func StartUI() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "scanme.nmap.org"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput:          ti,
		err:                nil,
		addressInputActive: true,
		progress:           progress.New(progress.WithDefaultGradient()),
		scanned:            0,
		openPorts:          make([]uint, 0),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			m.addressInputActive = false
			return m, startScanning(m.textInput.Value())
		}

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	case results:
		m.openPorts = msg.ports
		m.scanned = msg.scanned

		if msg.finished {
			cmd = m.progress.SetPercent(100)
			return m, tea.Quit
		} else {
			percents := float64(m.scanned) / float64(Ports)
			cmd = m.progress.SetPercent(percents)
		}

		return m, tea.Batch(cmd, tickProgressBar())

	case errMsg:
		m.err = msg
		return m, nil
	}

	if m.addressInputActive {
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	return m, tickProgressBar()
}

func (m model) View() string {
	inputView := fmt.Sprintf(
		"Enter address for port scaning\n\n%s\n\n%s",
		m.textInput.View(),
		helpStyle("(Ctrl+C to quit)\n\n"),
	)
	if m.addressInputActive {
		return inputView
	}

	pad := strings.Repeat(" ", padding)

	portsView := "Open ports: \n"
	for _, v := range m.openPorts {
		portsView += fmt.Sprintf("\t%d\n", v)
	}
	return inputView +
		pad + m.progress.View() + "\n\n" + portsView
}

func startScanning(address string) tea.Cmd {
	StartScanning(address)
	return tickProgressBar()
}

func tickProgressBar() tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		finished, foundPorts, scanned := GetFoundPorts()
		return results{
			ports:    foundPorts,
			finished: finished,
			scanned:  scanned,
		}
	})
}
