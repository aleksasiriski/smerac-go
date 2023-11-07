package calendar

import (
	"time"

	"google.golang.org/api/calendar/v3"
)

type Weekday struct {
	Name  string
	Items []*calendar.Event
}

type Info struct {
	Name  string
	Start []time.Time
	End   []time.Time
}

type Week struct {
	Mon Weekday
	Tue Weekday
	Wed Weekday
	Thu Weekday
	Fri Weekday
	Sat Weekday
	Sun Weekday
}

type ItemParsed struct {
	Name  string
	Infos []Info
}

type WeekdayParsed struct {
	Name  string
	Items []ItemParsed
}

type WeekParsed struct {
	Mon WeekdayParsed
	Tue WeekdayParsed
	Wed WeekdayParsed
	Thu WeekdayParsed
	Fri WeekdayParsed
	Sat WeekdayParsed
	Sun WeekdayParsed
}

type WeekOutput struct {
	Mon string
	Tue string
	Wed string
	Thu string
	Fri string
	Sat string
	Sun string
}
