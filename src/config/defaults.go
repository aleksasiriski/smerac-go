package config

func New() *Config {
	return &Config{
		Google:    Google{},
		Calendars: []Calendar{},
		Days: NamedDays{
			Monday:    "Monday",
			Tuesday:   "Tuesday",
			Wednesday: "Wednesday",
			Thursday:  "Thursday",
			Friday:    "Friday",
			Saturday:  "Saturday",
			Sunday:    "Sunday",
		},
	}
}
