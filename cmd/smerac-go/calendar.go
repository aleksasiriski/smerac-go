package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sourcegraph/conc"

	"github.com/bwmarrin/discordgo"
	"google.golang.org/api/calendar/v3"
)

type ItemTime struct {
	DateTime time.Time `json:"dateTime"`
}

type Item struct {
	Id      string    `json:"id"`
	Summary string    `json:"summary"`
	Updated time.Time `json:"updated"`
	Start   ItemTime  `json:"start"`
	End     ItemTime  `json:"end"`
}

type CalendarJson struct {
	Summary string    `json:"summary"`
	Updated time.Time `json:"updated"`
	Items   []Item    `json:"items"`
}

type Weekday struct {
	Name  string
	Items []Item
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

func (day *Weekday) Sort() {
	sort.Slice(day.Items, func(i, j int) bool {
		return day.Items[i].Start.DateTime.Before(day.Items[j].Start.DateTime)
	})
}

func (week *Week) Sort() {
	var worker conc.WaitGroup

	worker.Go(func() {
		week.Mon.Sort()
	})
	worker.Go(func() {
		week.Tue.Sort()
	})
	worker.Go(func() {
		week.Wed.Sort()
	})
	worker.Go(func() {
		week.Thu.Sort()
	})
	worker.Go(func() {
		week.Fri.Sort()
	})
	worker.Go(func() {
		week.Sat.Sort()
	})
	worker.Go(func() {
		week.Sun.Sort()
	})

	worker.Wait()
}

func (week *Week) Generate(items []Item) {
	foundItems := false

	for _, item := range items {
		currentTime := time.Now()
		if item.Start.DateTime.After(currentTime) && item.Start.DateTime.Before(currentTime.AddDate(0, 0, 7)) {
			log.Trace().
				Str("item", fmt.Sprintf("%v", item)).
				Msg("Appending")
			switch item.Start.DateTime.Weekday().String() {
			case "Monday":
				{
					week.Mon.Items = append(week.Mon.Items, item)
				}
			case "Tuesday":
				{
					week.Tue.Items = append(week.Tue.Items, item)
				}
			case "Wednesday":
				{
					week.Wed.Items = append(week.Wed.Items, item)
				}
			case "Thursday":
				{
					week.Thu.Items = append(week.Thu.Items, item)
				}
			case "Friday":
				{
					week.Fri.Items = append(week.Fri.Items, item)
				}
			case "Saturday":
				{
					week.Sat.Items = append(week.Sat.Items, item)
				}
			case "Sunday":
				{
					week.Sun.Items = append(week.Sun.Items, item)
				}
			}
			foundItems = true
		} else {
			log.Trace().
				Str("item start", fmt.Sprintf("%v", item.Start.DateTime)).
				Str("item end", fmt.Sprintf("%v", item.End.DateTime)).
				Str("current time", fmt.Sprintf("%v", currentTime)).
				Msg("Item not within time scope")
		}
	}

	if !foundItems {
		log.Error().
			Msg("No items found")
	}

	week.Sort()
}

func (day *Weekday) Parse() WeekdayParsed {
	dayParsed := WeekdayParsed{Name: day.Name, Items: make([]ItemParsed, 0)}

	for _, item := range day.Items {
		nameinfo := strings.SplitN(item.Summary, ",", 2)
		name := nameinfo[0]
		info := nameinfo[1]
		start := item.Start.DateTime
		end := item.End.DateTime

		foundEvent := false
		for _, itemParsed := range dayParsed.Items {
			compName := strings.ReplaceAll(strings.ToUpper(name), " ", "")
			compNameParsed := strings.ReplaceAll(strings.ToUpper(itemParsed.Name), " ", "")

			if compName == compNameParsed {
				foundInfo := false
				for _, infoParsed := range itemParsed.Infos {
					compInfo := strings.ReplaceAll(strings.ToUpper(info), " ", "")
					compInfoParsed := strings.ReplaceAll(strings.ToUpper(infoParsed.Name), " ", "")

					if compInfo == compInfoParsed {
						infoParsed.Start = append(infoParsed.Start, start)
						infoParsed.End = append(infoParsed.End, end)

						foundInfo = true
						break
					}
				}

				if !foundInfo {
					newInfo := Info{
						Name:  info,
						Start: make([]time.Time, 1),
						End:   make([]time.Time, 1),
					}
					newInfo.Start[0] = start
					newInfo.End[0] = end

					itemParsed.Infos = append(itemParsed.Infos, newInfo)
				}

				foundEvent = true
				break
			}
		}

		if !foundEvent {
			newInfo := Info{
				Name:  info,
				Start: make([]time.Time, 1),
				End:   make([]time.Time, 1),
			}
			newInfo.Start[0] = start
			newInfo.End[0] = end

			newItemParsed := ItemParsed{
				Name:  name,
				Infos: make([]Info, 1),
			}
			newItemParsed.Infos[0] = newInfo

			dayParsed.Items = append(dayParsed.Items, newItemParsed)
		}
	}

	return dayParsed
}

func (week *Week) Parse() WeekParsed {
	weekParsed := WeekParsed{}
	var worker conc.WaitGroup

	worker.Go(func() {
		weekParsed.Mon = week.Mon.Parse()
	})
	worker.Go(func() {
		weekParsed.Tue = week.Tue.Parse()
	})
	worker.Go(func() {
		weekParsed.Wed = week.Wed.Parse()
	})
	worker.Go(func() {
		weekParsed.Thu = week.Thu.Parse()
	})
	worker.Go(func() {
		weekParsed.Fri = week.Fri.Parse()
	})
	worker.Go(func() {
		weekParsed.Sat = week.Sat.Parse()
	})
	worker.Go(func() {
		weekParsed.Sun = week.Sun.Parse()
	})

	log.Trace().
		Msg("Waiting for parse")
	worker.Wait()
	return weekParsed
}

func (day *WeekdayParsed) Stringify() string {
	spacer := "-------------------------"
	output := spacer + "\n\n**" + day.Name + ":**\n\n"

	for _, item := range day.Items {
		output += "--- " + item.Name + " ---\n"

		for _, info := range item.Infos {
			output += info.Name + "\n"

			for index := range info.Start {
				output += "**" + info.Start[index].Format("03:04") + "** - " + info.End[index].Format("03:04") + "\n"
			}
		}

		output += "\n"
	}
	output += spacer + "\n"

	return output
}

func generateAndParseWeek(items []Item) WeekParsed {
	week := Week{
		Mon: Weekday{
			Name: "Monday",
		},
		Tue: Weekday{
			Name: "Tuesday",
		},
		Wed: Weekday{
			Name: "Wednesday",
		},
		Thu: Weekday{
			Name: "Thursday",
		},
		Fri: Weekday{
			Name: "Friday",
		},
		Sat: Weekday{
			Name: "Saturday",
		},
		Sun: Weekday{
			Name: "Sunday",
		},
	}

	week.Generate(items)
	weekParsed := week.Parse()
	log.Debug().
		Str("week", fmt.Sprintf("%v", weekParsed)).
		Msg("Parsed")

	return weekParsed
}

func updateCalendar(calendar Calendar, discord *discordgo.Session, google Google) (WeekParsed, error) {
	week := WeekParsed{}

	url := fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/%s@group.calendar.google.com/events?key=%s", calendar.Id, google.Token)

	log.Trace().
		Str("url", url).
		Msg("Getting calendar via API")

	ctx := context.Background()
	calendarService, err := calendar.NewService(ctx, option.WithAPIKey(google.Token))

	resp, err := http.Get(url)
	if err != nil {
		return week, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return week, fmt.Errorf("request was not OK: %s", resp.Status)
	}

	dec := json.NewDecoder(resp.Body)
	var calendarJson CalendarJson
	err = dec.Decode(&calendarJson)
	if err != nil {
		return week, err
	}
	log.Debug().
		Str("Calendar JSON", fmt.Sprintf("%v", calendarJson)).
		Msg("Decoded API response")

	week = generateAndParseWeek(calendarJson.Items)
	return week, nil
}

func outputWeek(week WeekParsed, discord *discordgo.Session) error {
	weekOutput := WeekOutput{}

	var worker conc.WaitGroup

	worker.Go(func() {
		weekOutput.Mon = week.Mon.Stringify()
	})
	worker.Go(func() {
		weekOutput.Tue = week.Tue.Stringify()
	})
	worker.Go(func() {
		weekOutput.Wed = week.Wed.Stringify()
	})
	worker.Go(func() {
		weekOutput.Thu = week.Thu.Stringify()
	})
	worker.Go(func() {
		weekOutput.Fri = week.Fri.Stringify()
	})
	worker.Go(func() {
		weekOutput.Sat = week.Sat.Stringify()
	})
	worker.Go(func() {
		weekOutput.Sun = week.Sun.Stringify()
	})

	log.Trace().
		Msg("Waiting for stringification")
	worker.Wait()

	//! temporary
	fmt.Print(weekOutput.Mon)
	fmt.Print(weekOutput.Tue)
	fmt.Print(weekOutput.Wed)

	return nil
}

func updateCalendars(calendars []Calendar, discord *discordgo.Session, google Google) error {
	var worker conc.WaitGroup
	for _, calendarIterator := range calendars {
		calendar := calendarIterator
		worker.Go(func() {
			for {
				log.Debug().
					Str("name", calendar.Name).
					Msg("Updating calendar")

				week, err := updateCalendar(calendar, discord, google)
				if err != nil {
					log.Error().
						Err(err).
						Msg(fmt.Sprintf("Failed while updating calendar %s:", calendar.Name))
				}

				log.Trace().
					Str("name", calendar.Name).
					Msg("Outputting calendar")

				err = outputWeek(week, discord)
				if err != nil {
					log.Error().
						Err(err).
						Msg(fmt.Sprintf("Failed while outputting calendar %s:", calendar.Name))
				}

				log.Trace().
					Str("name", calendar.Name).
					Msg("Sleeping calendar")

				if calendar.TimeBetweenChecks == 0 {
					time.Sleep(time.Hour * 3)
				} else {
					time.Sleep(time.Hour * time.Duration(calendar.TimeBetweenChecks))
				}
			}
		})
	}
	return nil
}
