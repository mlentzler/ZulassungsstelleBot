package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mlentzler/ZulassungsstelleBot/internal/domain"
)

func LoadMenu(path string) (domain.MenuNode, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return domain.MenuNode{}, fmt.Errorf("menu read: %w", err)
	}
	var root domain.MenuNode
	if err := json.Unmarshal(b, &root); err != nil {
		return domain.MenuNode{}, fmt.Errorf("menu parse: %w", err)
	}
	return root, nil
}
