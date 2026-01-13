package main

import (
	"log"
	"time"
	"govoid/processlist"
	"govoid/helper"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/*
This example assumes an existing understanding of commands and messages. If you
haven't already read our tutorials on the basics of Bubble Tea and working
with commands, we recommend reading those first.

Find them at:
https://github.com/charmbracelet/bubbletea/tree/master/tutorials/commands
https://github.com/charmbracelet/bubbletea/tree/master/tutorials/basics
*/

// sessionState is used to track which model is focused
type sessionState uint

const (
	defaultTime              = time.Minute
	timerView   sessionState = iota
	spinnerView
)

var (
	// Available spinners
	spinners = []spinner.Spinner{
		spinner.Line,
		spinner.Dot,
		spinner.MiniDot,
		spinner.Jump,
		spinner.Pulse,
		spinner.Points,
		spinner.Globe,
		spinner.Moon,
		spinner.Monkey,
	}
	appList = lipgloss.NewStyle().
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
	deviceStats = lipgloss.NewStyle().
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))

	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type mainModel struct {
	state   sessionState
	timer   timer.Model
	spinner spinner.Model
	index   int
	width   int
	height  int

	processList processlist.Model
}

func newModel(timeout time.Duration) mainModel {
	m := mainModel{state: timerView}
	m.timer = timer.New(timeout)
	m.spinner = spinner.New()

	data := helper.GetApps()
	helper.ReverseSlice(data)
    m.processList = processlist.New(data, 0, 0)
	return m
}

func (m mainModel) Init() tea.Cmd {
	// start the timer and spinner on program start
	return tea.Batch(m.timer.Init(), m.spinner.Tick)
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.processList, cmd = m.processList.Update(msg)
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        listWidth := (m.width * 70 / 100) - 2
        listHeight := m.height - 5
        
        m.processList.List.SetWidth(listWidth)
        m.processList.List.SetHeight(listHeight)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.state == timerView {
				m.state = spinnerView
			} else {
				m.state = timerView
			}
		case "n":
			if m.state == timerView {
				m.timer = timer.New(defaultTime)
				cmds = append(cmds, m.timer.Init())
			} else {
				m.Next()
				m.resetSpinner()
				cmds = append(cmds, m.spinner.Tick)
			}
		case "p":
			if m.state == timerView {
				// kill the selected app
				selectedItem := m.processList.List.SelectedItem()
				log.Printf("Selected item to kill: %v", selectedItem)
				if selectedItem != nil {
					appName := string(selectedItem.(processlist.Item))
					err := helper.KillApp(appName)
					if err != nil {
						//log.Printf("Failed to kill app %s: %v", appName, err)
					} else {
						// Refresh the process list after killing the app
						data := helper.GetApps()
						helper.ReverseSlice(data)
						m.processList = processlist.New(data, m.processList.List.Width(), m.processList.List.Height())
					}
				}
			}
		}
		switch m.state {
		// update whichever model is focused
		case spinnerView:
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		default:
			m.timer, cmd = m.timer.Update(msg)
			cmds = append(cmds, cmd)
		}
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	case timer.TickMsg:
		m.timer, cmd = m.timer.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}
func (m mainModel) View() string {
    if m.width == 0 {
        return "loading..."
    }

    listWidth := (m.width * 70 / 100) - 2
    statsWidth := (m.width * 30 / 100) - 2
    
    // Apply dimensions to the styles
    appList = appList.Width(listWidth).Height(m.height - 5)
    deviceStats = deviceStats.Width(statsWidth).Height(m.height - 5)

    // Highlight the active border
    if m.state == timerView {
        appList = appList.BorderForeground(lipgloss.Color("205")) // Pink focus
        deviceStats = deviceStats.BorderForeground(lipgloss.Color("69"))
    } else {
        appList = appList.BorderForeground(lipgloss.Color("69"))
        deviceStats = deviceStats.BorderForeground(lipgloss.Color("205")) // Pink focus
    }

    s := lipgloss.JoinHorizontal(lipgloss.Top, 
        appList.Render(m.processList.View()), 
        deviceStats.Render(m.spinner.View()),
    )
    
    s += helpStyle.Render("\ntab: switch focus • q: exit • p: kill\n")
    return s
}

func (m mainModel) currentFocusedModel() string {
	if m.state == timerView {
		return "timer"
	}
	return "spinner"
}

func (m *mainModel) Next() {
	if m.index == len(spinners)-1 {
		m.index = 0
	} else {
		m.index++
	}
}

func (m *mainModel) resetSpinner() {
	m.spinner = spinner.New()
	m.spinner.Style = spinnerStyle
	m.spinner.Spinner = spinners[m.index]
}

func main() {
	p := tea.NewProgram(newModel(defaultTime), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}