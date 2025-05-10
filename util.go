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

func LoadConfig() (map[string]any, error) {
	var cfg map[string]any

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

	_, ok := cfg["general"].(map[string]any)
	if !ok {
		import_array := map[string]any{"import": []any{}}
		cfg["general"] = import_array
		log.Printf("no general settings found... creating one: %v", cfg)
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

	general, ok := config["general"].(map[string]any)
	if !ok {
		log.Fatalf("failed to load general settings")
	}

	import_array, ok := general["import"].([]any)
	if ok {
		data := fmt.Sprintf("~/.config/alacritty/themes/%s.toml", th.name)
		import_array = import_array[:max(0, len(import_array)-2)]
		import_array = append(import_array, data)

		general["import"] = import_array
		config["general"] = general
	} else {
		log.Println(general)
		log.Fatalf("couldnt parse config file")
	}

	b, err := toml.Marshal(config)
	if err != nil {
		log.Panic(err)
	}

	saveConfig(b)
}
