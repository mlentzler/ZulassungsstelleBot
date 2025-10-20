package chromedpdrv

var XpCookieTry = []string{
	`//*[self::button or self::a][normalize-space(.)="OK"]`,
	`//*[self::button or self::a][contains(normalize-space(.),"Akzeptier")]`,
	`//*[self::button or self::a][contains(normalize-space(.),"Zustimm")]`,
	`//*[self::button or self::a][contains(normalize-space(.),"Einverstanden")]`,
}

// „Termin buchen“-Button
const XpTerminBuchen = `//*[self::a or self::button][contains(normalize-space(.),"Termin buchen")]`

// Zeit-Buttons (HH:MM)
const XpTimeButtons = `//button[normalize-space(.) and not(@disabled)] | //a[normalize-space(.) and not(@disabled)]`
