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

// Formularfelder (neue Struktur)
const (
	XpInputName         = `//span[contains(., 'Nachname, Vorname')]/ancestor::label//input`
	XpInputEmail        = `//span[contains(., 'E-Mail-Adresse')]/ancestor::label//input`
	XpInputPhone        = `//span[contains(., 'Handy/Telefon')]/ancestor::label//input`
	XpCheckboxPrivacy   = `label[for="IsTermsOfServiceConsentObtained"]` // Target the label for the checkbox
)

// Submit Button
const XpSubmit = `//button[contains(normalize-space(.), 'Weiter')]` // Find by text
