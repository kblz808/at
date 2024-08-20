package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const defaultWidth = 20

func main() {
	items, err := getThemes()
	if err != nil {
		log.Panic(err)
	}

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
