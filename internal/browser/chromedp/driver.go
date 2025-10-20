package chromedpdrv

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"

	"github.com/mlentzler/ZulassungsstelleBot/internal/browser"
)

type Driver struct {
	sess *Session
	loc  *time.Location
}

func NewDriver(headless bool, loc *time.Location) (*Driver, error) {
	s, err := New(headless)
	if err != nil {
		return nil, err
	}
	return &Driver{sess: s, loc: loc}, nil
}

func (d *Driver) Open(ctx context.Context) error  { return nil }
func (d *Driver) Close(ctx context.Context) error { d.sess.Close(); return nil }

func (d *Driver) StartFlow(ctx context.Context, baseURL string, titles []string, selectors []string) error {
	c := d.sess.Context()

	// Start
	if err := chromedp.Run(c,
		chromedp.Navigate(baseURL),
		chromedp.WaitReady("body", chromedp.ByQuery),
	); err != nil {
		return fmt.Errorf("navigate: %w", err)
	}

	// Menü-Selektoren nacheinander klicken
	for i := range selectors {
		sel := selectors[i]
		title := ""
		if i < len(titles) {
			title = titles[i]
		}

		stepCtx, cancel := context.WithTimeout(c, 8*time.Second)
		defer cancel()

		var err error
		if strings.HasPrefix(sel, "/") || strings.HasPrefix(sel, "(") {
			err = chromedp.Run(stepCtx,
				chromedp.WaitVisible(sel, chromedp.BySearch),
				chromedp.ScrollIntoView(sel, chromedp.BySearch),
				chromedp.Click(sel, chromedp.NodeVisible, chromedp.BySearch),
			)
		} else {
			err = chromedp.Run(stepCtx,
				chromedp.WaitVisible(sel, chromedp.ByQuery),
				chromedp.ScrollIntoView(sel, chromedp.ByQuery),
				chromedp.Click(sel, chromedp.NodeVisible, chromedp.ByQuery),
			)
		}
		if err != nil {
			return fmt.Errorf("menu step %d failed (title=%q sel=%q): %w", i+1, title, sel, err)
		}
		_ = chromedp.Run(c, Sleep(400))
	}

	// Produktseite → “Termin buchen”
	return chromedp.Run(c,
		chromedp.WaitVisible(XpTerminBuchen, chromedp.BySearch),
		chromedp.ScrollIntoView(XpTerminBuchen, chromedp.BySearch),
		chromedp.Click(XpTerminBuchen, chromedp.NodeVisible, chromedp.BySearch),
		Sleep(500),
	)
}

func (d *Driver) PickDate(ctx context.Context, date time.Time) error {
	return nil
}

var (
	reTime = regexp.MustCompile(`^\s*\d{1,2}:\d{2}\s*$`)

	xpDateGroups = []string{
		`//*[@data-date]`,
		`//section[contains(@class,"day") or contains(@class,"date")]`,
		`//div[contains(@class,"day") or contains(@class,"date")]`,
	}
	xpTimesWithin = `.//button[normalize-space(.) and not(@disabled)] | .//a[normalize-space(.) and not(@disabled)]`
)

func parseDateLabel(s string, loc *time.Location) (time.Time, bool) {
	s = strings.TrimSpace(s)
	// YYYY-MM-DD
	if t, err := time.ParseInLocation("2006-01-02", s, loc); err == nil {
		return t, true
	}
	// DD.MM.YYYY
	if t, err := time.ParseInLocation("02.01.2006", s, loc); err == nil {
		return t, true
	}
	// "Montag, 21.10.2025" → nach Komma
	if i := strings.LastIndex(s, ","); i >= 0 {
		if t, err := time.ParseInLocation("02.01.2006", strings.TrimSpace(s[i+1:]), loc); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func (d *Driver) ListSlots(ctx context.Context) ([]browser.Slot, error) {
	c := d.sess.Context()

	// Alle Slot-Links global holen: <a class="... time-container ..." onclick="selectTime(..., '2025-10-23T08:05:00+02:00', ...)"
	const xpAllTimeLinks = `//a[contains(@class,"time-container") and contains(@onclick,"selectTime")]`

	var nodes []*cdp.Node
	if err := chromedp.Run(c, chromedp.Nodes(xpAllTimeLinks, &nodes, chromedp.BySearch)); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, nil
	}

	// ISO-Zeitstempel aus dem onclick ziehen
	reISO := regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:Z|[+-]\d{2}:\d{2})`)

	out := make([]browser.Slot, 0, len(nodes))
	for _, n := range nodes {
		onclick, _ := getAttr(n, "onclick") // ← direkt aus dem Node, NICHT via XPath
		if onclick == "" {
			continue
		}
		iso := reISO.FindString(onclick)
		if iso == "" {
			// notfalls aria-label prüfen, aber onclick ist der Goldstandard
			continue
		}
		t, err := time.Parse(time.RFC3339, iso)
		if err != nil {
			continue
		}
		out = append(out, browser.Slot{
			Start: t.In(d.loc),
			Ref:   n, // Node direkt für MouseClickNode
		})
	}

	return out, nil
}

// getAttr liest ein Attribut direkt aus dem CDP-Node (Attribute-Liste: [name, value, name, value, ...])
func getAttr(n *cdp.Node, key string) (string, bool) {
	for i := 0; i+1 < len(n.Attributes); i += 2 {
		if n.Attributes[i] == key {
			return n.Attributes[i+1], true
		}
	}
	return "", false
}

func (d *Driver) BookSlot(ctx context.Context, s browser.Slot, form map[string]string) error {
	c := d.sess.Context()
	n, _ := s.Ref.(*cdp.Node)

	// 1) aria-label vom Node lesen und daraus einen exakten XPath bauen
	aria, _ := getAttr(n, "aria-label")
	var actions []chromedp.Action
	if aria != "" {
		xp := `//a[contains(@class,"time-container") and @aria-label=` + xpathQuote(aria) + `]`
		actions = append(actions,
			chromedp.WaitVisible(xp, chromedp.BySearch),
			chromedp.ScrollIntoView(xp, chromedp.BySearch),
			chromedp.Click(xp, chromedp.NodeVisible, chromedp.BySearch),
		)
	} else {
		// Fallback: direkter Node-Klick (falls aria-label fehlt)
		actions = append(actions,
			chromedp.MouseClickNode(n),
		)
	}
	actions = append(actions, Sleep(350))

	// 2) Formular abwarten – ein typisches Label/Submit sichtbar?
	actions = append(actions,
		// passe Labeltexte an, wenn nötig
		chromedp.WaitReady(`//label[contains(normalize-space(.),"Name")]`, chromedp.BySearch),
	)

	// 3) Felder füllen
	if v := form["name"]; v != "" {
		actions = append(actions, chromedp.SetValue(XpInputByAnyLabel("Name", "Vor- und Nachname"), v, chromedp.BySearch))
	}
	if v := form["email"]; v != "" {
		actions = append(actions, chromedp.SetValue(XpInputByAnyLabel("E-Mail", "E-Mail-Adresse"), v, chromedp.BySearch))
	}
	if v := form["telefon"]; v != "" {
		actions = append(actions, chromedp.SetValue(XpInputByAnyLabel("Telefon", "Telefonnummer"), v, chromedp.BySearch))
	}

	// 4) Bestätigen
	actions = append(actions,
		Sleep(250),
		chromedp.Click(XpSubmit, chromedp.NodeVisible, chromedp.BySearch),
		chromedp.WaitReady("body", chromedp.ByQuery),
	)

	return chromedp.Run(c, actions...)
}

func xpathQuote(s string) string { return `"` + s + `"` }

func nodeXPath(n *cdp.Node) string {
	return fmt.Sprintf(`//*[@node-id="%d"]`, n.NodeID)
}

func XpInputByAnyLabel(labels ...string) string {
	parts := make([]string, 0, len(labels))
	for _, l := range labels {
		parts = append(parts, `(//label[contains(normalize-space(.), "`+l+`")]/following::input[1])[1]`)
	}
	return strings.Join(parts, " | ")
}

const XpSubmit = `//*[self::button or self::a][
  contains(translate(normalize-space(.),"ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÜ","abcdefghijklmnopqrstuvwxyzäöü"),"bestät")
  or contains(normalize-space(.),"Buchen")
]`
