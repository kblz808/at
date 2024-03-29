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

	return cfg, nil
}

func saveConfig(config []byte) error {
	err := os.WriteFile(filePath, config, 0644)
	return err
}

func applyTheme(th theme) {
	config, err := LoadConfig()
	if err != nil {
		log.Panic(err)
	}

	import_array, ok := config["import"].([]interface{})
	if ok {
		data := fmt.Sprintf("~/.config/alacritty/themes/themes/%s.toml", th.name)
		import_array = import_array[:len(import_array)-1]
		import_array = append(import_array, data)
		config["import"] = import_array
	} else {
		println("not ok")
	}

	b, err := toml.Marshal(config)
	if err != nil {
		log.Panic(err)
	}

	saveConfig(b)
}
