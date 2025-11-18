package chromedpdrv

var XpCookieTry = []string{
	`//*[self::button or self::a][normalize-space(.)="OK"]`,
	`//*[self::button or self::a][contains(normalize-space(.),"Akzeptier")]`,
	`//*[self::button or self::a][contains(normalize-space(.),"Zustimm")]`,
	`//*[self::button or self::a][contains(normalize-space(.),"Einverstanden")]`,
}

const XpBookSloot = `//*[self::a or self::button][contains(normalize-space(.),"Termin buchen")]`

const XpTimeButtons = `//button[normalize-space(.) and not(@disabled)] | //a[normalize-space(.) and not(@disabled)]`

const (
	XpInputName       = `//span[contains(., 'Nachname, Vorname')]/ancestor::label//input`
	XpInputEmail      = `//span[contains(., 'E-Mail-Adresse')]/ancestor::label//input`
	XpInputPhone      = `//span[contains(., 'Handy/Telefon')]/ancestor::label//input`
	XpCheckboxPrivacy = `label[for="IsTermsOfServiceConsentObtained"]`
)

const (
	XpContinue       = `//button[contains(normalize-space(.), 'Weiter')]`
	XpConfirmBooking = `//button[contains(normalize-space(.), 'Bestätigen')]`
)
