package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

var filePath string

func getThemes() ([]list.Item, error) {
	items := []list.Item{}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}

	folderPath := filepath.Join(homeDir, ".config", "alacritty", "themes")

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
