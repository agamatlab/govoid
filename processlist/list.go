package processlist

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	ItemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type Item string

func (i Item) FilterValue() string { return "" }

type ItemDelegate struct{}

func (d ItemDelegate) Height() int                             { return 1 }
func (d ItemDelegate) Spacing() int                            { return 0 }
func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type Model struct {
	List     list.Model
	choice   string
	quitting bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func New(data []string, width, height int) Model {
    Items := []list.Item{}
    for _, val := range data {
        Items = append(Items, Item(val))
    }

    l := list.New(Items, ItemDelegate{}, width, height)
    l.Title = "Active Processes"
    l.SetShowStatusBar(false)
    l.SetFilteringEnabled(false)
    l.Styles.Title = titleStyle
    l.Styles.PaginationStyle = paginationStyle
    l.Styles.HelpStyle = helpStyle

    return Model{List: l}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "enter" {
            i, ok := m.List.SelectedItem().(Item)
            if ok {
                m.choice = string(i)
            }
        }
    }
    
    var cmd tea.Cmd
    m.List, cmd = m.List.Update(msg)
    return m, cmd
}

func (m Model) View() string {
    return m.List.View()
}

