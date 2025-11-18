package chromedpdrv

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
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

type slotRef struct {
	Node *cdp.Node
	ISO  string
	Aria string
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

	if err := chromedp.Run(c,
		chromedp.Navigate(baseURL),
		chromedp.WaitReady("body", chromedp.ByQuery),
	); err != nil {
		return fmt.Errorf("navigate: %w", err)
	}

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

	return chromedp.Run(c,
		chromedp.WaitVisible(XpBookSloot, chromedp.BySearch),
		chromedp.ScrollIntoView(XpBookSloot, chromedp.BySearch),
		chromedp.Click(XpBookSloot, chromedp.NodeVisible, chromedp.BySearch),
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

// robustness against different date formats
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

	const xpCandidates = `
(
  //a[contains(@class,"time") or contains(@class,"slot") or contains(@class,"time-container") or contains(@onclick,"selectTime") or @aria-label]
 |//button[contains(@class,"time") or contains(@class,"slot") or contains(@onclick,"selectTime") or @aria-label]
)
[not(@disabled)]
`

	var nodes []*cdp.Node
	if err := chromedp.Run(c, chromedp.Nodes(xpCandidates, &nodes, chromedp.BySearch)); err != nil {
		return nil, err
	}

	if len(nodes) == 0 {
		return nil, nil
	}

	reISO := regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:Z|[+-]\d{2}:\d{2})`)
	reTime := regexp.MustCompile(`\b(\d{1,2}):(\d{2})\b`)

	out := make([]browser.Slot, 0, len(nodes))
	for _, n := range nodes {
		onclick, _ := getAttr(n, "onclick")
		aria, _ := getAttr(n, "aria-label")
		dataDatetime, _ := getAttr(n, "data-datetime")
		dataDate, _ := getAttr(n, "data-date")

		var (
			iso string
			ts  time.Time
			ok  bool
		)

		if onclick != "" {
			if m := reISO.FindString(onclick); m != "" {
				iso = m
			}
		}
		if iso == "" && dataDatetime != "" && reISO.MatchString(dataDatetime) {
			iso = dataDatetime
		}

		if iso != "" {
			t, err := time.Parse(time.RFC3339, iso)
			if err == nil {
				ts, ok = t.In(d.loc), true
			} else {
				d.logf("ListSlots: ISO %q parse error: %v", iso, err)
			}
		}

		// Fallback: use aria-label to parse time
		if !ok && aria != "" {
			if day, dayOK := parseDateLabel(aria, d.loc); dayOK {
				if tm := reTime.FindString(aria); tm != "" {
					parts := strings.SplitN(tm, ":", 2)
					hh := parts[0]
					mm := parts[1]
					layout := "2006-01-02 15:04"
					composed := fmt.Sprintf("%04d-%02d-%02d %s:%s",
						day.Year(), int(day.Month()), day.Day(), hh, mm)
					if t, err := time.ParseInLocation(layout, composed, d.loc); err == nil {
						ts, ok = t, true
					} else {
						d.logf("ListSlots: aria-label compose parse error (%q): %v", composed, err)
					}
				}
			}
		}

		// 2nd Fallback: data-date + time from aria-label
		if !ok && dataDate != "" {
			if day, dayOK := parseDateLabel(dataDate, d.loc); dayOK {
				var hhmm string
				if aria != "" {
					hhmm = reTime.FindString(aria)
				}
				if hhmm != "" {
					layout := "2006-01-02 15:04"
					composed := fmt.Sprintf("%04d-%02d-%02d %s",
						day.Year(), int(day.Month()), day.Day(), hhmm)
					if t, err := time.ParseInLocation(layout, composed, d.loc); err == nil {
						ts, ok = t, true
					}
				}
			}
		}

		if !ok {
			d.logf("ListSlots: verwerfe Kandidat (kein ISO/Datum+Zeit) nodeID=%d onclick=%q aria=%q data-datetime=%q data-date=%q",
				n.NodeID, onclick, aria, dataDatetime, dataDate)
			continue
		}

		out = append(out, browser.Slot{
			Start: ts,
			Ref:   slotRef{Node: n, ISO: iso, Aria: aria},
		})
		d.logf("ListSlots: erkannt %s (ISO=%q, aria=%q)", ts.Format(time.RFC3339), iso, aria)
	}

	d.logf("ListSlots: insgesamt %d Slots erkannt", len(out))
	return out, nil
}

func (d *Driver) logf(msg string, args ...any) {
	if os.Getenv("DEBUG") == "true" {
		log.Printf("[chromedpdrv] "+msg, args...)
	}
}

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
	d.logf("BookSlot: AUFGERUFEN start=%s refType=%T form=%s", s.Start.Format(time.RFC3339), s.Ref, d.dumpFormMap(form))

	var (
		n    *cdp.Node
		iso  string
		aria string
	)
	switch r := s.Ref.(type) {
	case slotRef:
		n = r.Node
		iso = r.ISO
		aria = r.Aria
	case *cdp.Node:
		n = r
		if v, _ := getAttr(n, "aria-label"); v != "" {
			aria = v
		}
		if v, _ := getAttr(n, "data-datetime"); v != "" {
			iso = v
		}
		if iso == "" {
			if v, _ := getAttr(n, "onclick"); v != "" {
				if m := regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:Z|[+-]\d{2}:\d{2})`).FindString(v); m != "" {
					iso = m
				}
			}
		}
	default:
	}

	if v := form["_forceISO"]; v != "" {
		d.logf("BookSlot: FORCE ISO override aktiv: %s", v)
		iso = v
	}
	if v := form["_forceARIA"]; v != "" {
		d.logf("BookSlot: FORCE ARIA override aktiv: %s", v)
		aria = v
	}

	if iso == "" && !s.Start.IsZero() {
		iso = s.Start.Format(time.RFC3339)
		d.logf("BookSlot: kein ISO in Ref, benutze aus Start: %s", iso)
	}
	hhmm := s.Start.Format("15:04")

	d.logf("BookSlot: Klick-Ziel iso=%q aria=%q hhmm=%q nodePresent=%v", iso, aria, hhmm, n != nil)

	var xps []string
	if iso != "" {
		xps = append(xps,
			`//a[contains(@onclick,`+xpathQuote(iso)+`)]`,
			`//button[contains(@onclick,`+xpathQuote(iso)+`)]`,
			`//*[@data-datetime=`+xpathQuote(iso)+`]`,
		)
	}
	if aria != "" {
		xps = append(xps,
			`//*[@aria-label=`+xpathQuote(aria)+`]`,
		)
	}
	xps = append(xps,
		`//a[contains(@class,"time-container") and contains(normalize-space(.),`+xpathQuote(hhmm)+`)]`,
		`//button[contains(normalize-space(.),`+xpathQuote(hhmm)+`)]`,
	)

	// 1) XPath-Click-Versuche
	var clickErr error
	for i, xp := range xps {
		d.logf("BookSlot: Versuch %d XPath=%s", i+1, xp)
		stepCtx, cancel := context.WithTimeout(c, 3*time.Second)
		err := chromedp.Run(stepCtx,
			chromedp.WaitVisible(xp, chromedp.BySearch),
			chromedp.ScrollIntoView(xp, chromedp.BySearch),
			chromedp.Click(xp, chromedp.NodeVisible, chromedp.BySearch),
			Sleep(350),
		)
		cancel()
		if err == nil {
			d.logf("BookSlot: Klick erfolgreich (XPath %d)", i+1)
			clickErr = nil
			break
		}
		d.logf("BookSlot: Klick fehlgeschlagen (XPath %d): %v", i+1, err)
		clickErr = err
	}

	// 2) Fallback: direct Node-Klick
	if clickErr != nil && n != nil {
		d.logf("BookSlot: Fallback MouseClickNode auf nodeID=%d", n.NodeID)
		if err := chromedp.Run(c, chromedp.MouseClickNode(n), Sleep(300)); err != nil {
			d.logf("BookSlot: MouseClickNode fehlgeschlagen: %v", err)
			clickErr = err
		} else {
			clickErr = nil
		}
	}

	// 3) Fallback: JS-click
	if clickErr != nil {
		d.logf("BookSlot: JS-Fallback click() wird versucht…")
		var ok bool
		js := `
(function(){
  var sel = [];
  %s
  for (var i=0;i<sel.length;i++){
    var el = document.querySelector(sel[i]);
    if (el && !el.disabled) { el.click(); return true; }
  }
  return false;
})()`

		var parts []string
		if iso != "" {
			parts = append(parts, fmt.Sprintf(`sel.push('[data-datetime=%s]');`, xpathQuote(iso)))
			parts = append(parts, fmt.Sprintf(`sel.push('a[onclick*=%s]');`, xpathQuote(iso)))
			parts = append(parts, fmt.Sprintf(`sel.push('button[onclick*=%s]');`, xpathQuote(iso)))
		}
		if aria != "" {
			parts = append(parts, fmt.Sprintf(`sel.push('[aria-label=%s]');`, xpathQuote(aria)))
		}
		parts = append(parts, fmt.Sprintf(`sel.push('a.time-container'); sel.push('button');`))

		payload := fmt.Sprintf(js, strings.Join(parts, "\n  "))
		if err := chromedp.Run(c, chromedp.EvaluateAsDevTools(payload, &ok)); err != nil {
			d.logf("BookSlot: JS-Fallback Fehler: %v", err)
		} else {
			d.logf("BookSlot: JS-Fallback Ergebnis: %v", ok)
			if ok {
				clickErr = nil
			}
		}
	}

	if clickErr != nil {
		if n != nil {
			d.logf("BookSlot: gebe OuterHTML des Kandidaten aus (nodeID=%d)…", n.NodeID)
			xp := ""
			if aria != "" {
				xp = `//*[@aria-label=` + xpathQuote(aria) + `]`
			} else if iso != "" {
				xp = `//*[@data-datetime=` + xpathQuote(iso) + `]`
			}
			if xp != "" {
				var html string
				_ = chromedp.Run(c, chromedp.OuterHTML(xp, &html, chromedp.BySearch))
				if html != "" {
					d.logf("BookSlot: OuterHTML-Kandidat:\n%s", html)
				}
			}
		}
		return fmt.Errorf("BookSlot: kein Klick möglich: %w", clickErr)
	}

	d.logf("BookSlot: warte auf Formular/Step-2 Indikatoren…")
	indicators := []string{
		`//label[contains(translate(normalize-space(.),"ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÜ","abcdefghijklmnopqrstuvwxyzäöü"),"name")]`,
		`//label[contains(translate(normalize-space(.),"ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÜ","abcdefghijklmnopqrstuvwxyzäöü"),"e-mail")]`,
		`//button[contains(translate(normalize-space(.),"ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÜ","abcdefghijklmnopqrstuvwxyzäöü"),"weiter")]`,
		`//h2[contains(translate(normalize-space(.),"ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÜ","abcdefghijklmnopqrstuvwxyzäöü"),"angaben")]`,
	}

	if err := chromedp.Run(c, waitAnyVisible(indicators)); err != nil {
		d.logf("BookSlot: Formular-Indikatoren NICHT erschienen: %v", err)
	} else {
		d.logf("BookSlot: Formular-Indikator sichtbar.")
	}

	return nil
}

func (d *Driver) FillAndContinue(ctx context.Context, form map[string]string) error {
	d.logf("FillAndContinue: called")
	c := d.sess.Context()

	name := strings.TrimSpace(form["name"])
	email := strings.TrimSpace(form["email"])
	phone := strings.TrimSpace(form["telefon"])

	actions := []chromedp.Action{
		// Name
		chromedp.SendKeys(XpInputName, name, chromedp.BySearch),

		// E-Mail
		chromedp.SendKeys(XpInputEmail, email, chromedp.BySearch),

		// Telefon
		chromedp.SendKeys(XpInputPhone, phone, chromedp.BySearch),

		// Checkbox (AGB)
		chromedp.Click(XpCheckboxPrivacy, chromedp.ByQuery, chromedp.NodeVisible),

		// Submit („Weiter“)
		chromedp.WaitEnabled(XpContinue, chromedp.BySearch),
		chromedp.Click(XpContinue, chromedp.BySearch),
		chromedp.Sleep(1000 * time.Millisecond),
	}

	if err := chromedp.Run(c, actions...); err != nil {
		return fmt.Errorf("FillAndContinue: %w", err)
	}

	d.logf("FillAndContinue: Formular ausgefüllt und abgeschickt.")
	return nil
}

func (d *Driver) ConfirmBooking(ctx context.Context) error {
	d.logf("ConfirmBooking: called")
	c := d.sess.Context()

	actions := []chromedp.Action{
		chromedp.WaitEnabled(XpConfirmBooking, chromedp.BySearch),
		chromedp.Click(XpConfirmBooking, chromedp.BySearch),
		chromedp.Sleep(500 * time.Millisecond),
	}

	if err := chromedp.Run(c, actions...); err != nil {
		return fmt.Errorf("ConfirmBooking: %w", err)
	}

	d.logf("ConfirmBooking: Termin bestätigt.")
	return nil
}

func (d *Driver) dumpFormMap(form map[string]string) string {
	b, _ := json.Marshal(form)
	return string(b)
}

func waitAnyVisible(xps []string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		var lastErr error
		for _, xp := range xps {
			stepCtx, cancel := context.WithTimeout(ctx, 1500*time.Millisecond)
			lastErr = chromedp.Run(stepCtx,
				chromedp.WaitVisible(xp, chromedp.BySearch),
			)
			cancel()
			if lastErr == nil {
				return nil
			}
		}
		return lastErr
	}
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
