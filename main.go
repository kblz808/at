package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const defaultWidth = 20

func newList(items []list.Item) list.Model {
	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "select theme"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return l
}

func main() {
	items, err := getThemes()
	if err != nil {
		log.Fatalf(err.Error())
	}

	l := newList(items)
	m := model{list: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("error running the program", err)
		os.Exit(1)
	}
}
