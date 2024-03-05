package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

func getThemes() ([]list.Item, error) {
	items := []list.Item{}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}

	folderPath := filepath.Join(homeDir, ".config", "alacritty", "themes", "themes")

	files, err := os.ReadDir(folderPath)
	if err != nil {
		println(err.Error())
		return nil, err
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

	return items, nil

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
