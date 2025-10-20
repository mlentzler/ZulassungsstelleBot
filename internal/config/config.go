package config

import "os"

type Config struct {
	MenuPath string
	TZ       string
	BaseURL  string
	Headless bool
	PollMin  int
	PollMax  int
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
	base := os.Getenv("BASE_URL")
	if base == "" {
		base = "https://reservation.frontdesksuite.com/pinneberg/Termin/Home/Index?Culture=de&PageId=f3e3da57-3aeb-4f3c-8d22-bb44721210d5&ShouldStartReserveTimeFlow=False&ButtonId=00000000-0000-0000-0000-000000000000"
	}
	return Config{
		MenuPath: menu,
		TZ:       tz,
		BaseURL:  base,
		Headless: false,
		PollMin:  45,
		PollMax:  120,
	}
}
