package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 16

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type theme struct {
	name     string
	location string
}

func (t theme) FilterValue() string { return t.name }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(theme)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.name)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list     list.Model
	choice   theme
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(theme)
			if ok {
				m.choice = i
			}
			changeTheme(m.choice)
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return "\n" + m.list.View()
}

func changeTheme(th theme) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	filePath := filepath.Join(homeDir, ".config", "alacritty", "alacritty.toml")

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	lines := strings.Split(string(content), "\n")

	lines[1] = fmt.Sprintf("\t \"~/.config/alacritty/themes/themes/%s.toml\"", th.name)

	newContent := strings.Join(lines, "\n")

	err = os.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
}

func main() {
	items := []list.Item{}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	folderPath := filepath.Join(homeDir, ".config", "alacritty", "themes", "themes")

	files, err := os.ReadDir(folderPath)
	if err != nil {
		println(err.Error())
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))

		absolutePath := filepath.Join(folderPath, file.Name())

		items = append(items, theme{name: fileName, location: absolutePath})
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "select theme"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("error running the program", err)
		os.Exit(1)
	}
}
