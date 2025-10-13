package config

import "os"

type Config struct {
	MenuPath string // Pfad zu configs/menu.json
	TZ       string // "Europe/Berlin"
}

func Load() Config {
	menu := os.Getenv("MENU_PATH")
	if menu == "" {
		menu = "configs/menu.json"
	}
	tz := os.Getenv("TZ")
	if tz == "" {
		tz = "Europe/Berlin"
	}
	return Config{MenuPath: menu, TZ: tz}
}
