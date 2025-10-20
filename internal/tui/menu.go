package tui

import (
	"strings"

	"github.com/mlentzler/ZulassungsstelleBot/internal/domain"
)

func currentLevelNodes(m *Model) []domain.MenuNode {
	node := m.menuRoot
	for _, i := range m.menuStack {
		if i < 0 || i >= len(node.Children) {
			return nil
		}
		node = node.Children[i]
	}
	return node.Children
}

func currentPathTitles(m *Model) []string {
	var titles []string
	node := m.menuRoot
	for _, i := range m.menuStack {
		if i < 0 || i >= len(node.Children) {
			return titles
		}
		node = node.Children[i]
		titles = append(titles, node.Title)
	}
	return titles
}

func currentPathSelectors(m *Model, leafIdx int) []string {
	var sels []string
	node := m.menuRoot
	for _, idx := range m.menuStack {
		node = node.Children[idx]
		if node.Selector != "" {
			sels = append(sels, node.Selector)
		}
	}
	// Blatt
	children := currentLevelNodes(m)
	if leafIdx >= 0 && leafIdx < len(children) && children[leafIdx].Selector != "" {
		sels = append(sels, children[leafIdx].Selector)
	}

	return sels
}
func breadcrumb(m *Model) string {
	parts := []string{"Start"}
	parts = append(parts, currentPathTitles(m)...)
	return strings.Join(parts, " > ")
}

func pushMenu(m *Model, i int) {
	nodes := currentLevelNodes(m)
	if i < 0 || i >= len(nodes) {
		return
	}
	m.menuStack = append(m.menuStack, i)
	m.menuCursor = 0
}

func popMenu(m *Model) {
	if len(m.menuStack) == 0 {
		return
	}
	m.menuStack = m.menuStack[:len(m.menuStack)-1]
	m.menuCursor = 0
}

// updateMenu(m, msg) & viewMenu(m) bitte in model.go / view.go aufrufen
