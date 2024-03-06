package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/pelletier/go-toml/v2"
)

var filePath string

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

func LoadConfig() (map[string]interface{}, error) {
	var cfg map[string]interface{}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error: ", err)
		return nil, err
	}

	filePath = filepath.Join(homeDir, ".config", "alacritty", "alacritty.toml")

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("error: ", err)
		return nil, err
	}

	err = toml.Unmarshal(content, &cfg)
	if err != nil {
		log.Panic(err)
	}

	_, ok := cfg["import"]
	if !ok {
		log.Panic("import array not found")
	}

	// lines := strings.Split(string(content), "\n")

	// lines[1] = fmt.Sprintf("\t \"~/.config/alacritty/themes/themes/%s.toml\"", th.name)

	// newContent := strings.Join(lines, "\n")

	// err = os.WriteFile(filePath, []byte(newContent), 0644)
	// if err != nil {
	// fmt.Println("error: ", err)
	// return
	// }

	return cfg, nil
}

func saveConfig(config string) error {
	err := os.WriteFile(filePath, []byte(config), 0644)
	return err
}

func applyTheme(th theme) {
	config, err := LoadConfig()
	if err != nil {
		log.Panic(err)
	}

	import_array_interface, ok := config["import"].([]interface{})
	if ok {
		import_array_string := make([]string, len(import_array_interface))
		for i, v := range import_array_interface {
			switch value := v.(type) {
			case string:
				import_array_string[i] = value
			default:
				fmt.Printf("Failed to convert element at index %d to string\n", i)
			}
		}
		fmt.Printf("%+v\n", import_array_string)
	} else {
		println("not ok")
	}

	// saveConfig("")
}
