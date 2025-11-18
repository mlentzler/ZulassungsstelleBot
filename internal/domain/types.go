package domain

type AvailabilityKind string

const (
	AvailOneOff    AvailabilityKind = "oneoff"
	AvailRecurring AvailabilityKind = "recurring"
)

type OneOff struct {
	DateISO  string
	FromHour int
	ToHour   int
}

type DayWindow struct {
	Weekday  string
	FromHour int
	ToHour   int
}

type Recurring struct {
	Days []DayWindow
}

type Availability struct {
	Kind      AvailabilityKind
	OneOff    *OneOff
	Recurring *Recurring
}

type MenuNode struct {
	Title    string     `json:"title" yaml:"title"`
	Children []MenuNode `json:"children,omitempty" yaml:"children,omitempty"`
	Selector string     `json:"selector,omitempty"`
	Path     []string   `json:"path,omitempty"`
}

type MenuChoice struct {
	Path      []string
	Selectors []string
}

type BookingRequest struct {
	Name  string
	Email string
	Phone string
	Menu  MenuChoice
	Avail Availability
	TZ    string
}
