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

var configPath string

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

func LoadConfig() (map[string]interface{}, error) {
	var cfg map[string]interface{}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error: ", err)
		return nil, err
	}

	configPath = filepath.Join(homeDir, ".config", "alacritty", "alacritty.toml")

	content, err := os.ReadFile(configPath)
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
		cfg["import"] = []interface{}{""}
	}

	return cfg, nil
}

func saveConfig(config []byte) error {
	err := os.WriteFile(configPath, config, 0644)
	return err
}

func applyTheme(th theme) {
	config, err := LoadConfig()
	if err != nil {
		log.Panic(err)
	}

	import_array, ok := config["import"].([]interface{})
	if ok {
		data := fmt.Sprintf("~/.config/alacritty/themes/%s.toml", th.name)
		import_array = import_array[:len(import_array)-1]
		import_array = append(import_array, data)
		config["import"] = import_array
	} else {
		log.Fatalf("couldnt parse config file")
	}

	b, err := toml.Marshal(config)
	if err != nil {
		log.Panic(err)
	}

	saveConfig(b)
}
