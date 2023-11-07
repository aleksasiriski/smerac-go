package calendar

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sourcegraph/conc"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

	"github.com/aleksasiriski/smerac-go/src/config"
	"github.com/aleksasiriski/smerac-go/src/webhook"
)

func (week *Week) Generate(items []*calendar.Event) {
	foundItems := false

	for _, item := range items {
		log.Trace().
			Str("item", fmt.Sprintf("%v", item)).
			Msg("Appending")
		weekday, err := time.Parse(time.RFC3339, item.Start.DateTime)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed parsing time")
		} else {
			switch weekday.Weekday().String() {
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
		}
	}

	if !foundItems {
		log.Warn().
			Msg("No items found")
	}
}

func (day Weekday) Parse() WeekdayParsed {
	dayParsed := WeekdayParsed{
		Name:  day.Name,
		Items: make([]ItemParsed, 0),
	}

	for _, item := range day.Items {
		nameinfo := strings.SplitN(item.Summary, ",", 2)
		name := nameinfo[0]
		info := nameinfo[1]

		start, err := time.Parse(time.RFC3339, item.Start.DateTime)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed parsing time")
		}
		end, err := time.Parse(time.RFC3339, item.End.DateTime)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed parsing time")
		}

		foundEvent := false
		for itemIndex, itemParsed := range dayParsed.Items {
			compName := strings.ReplaceAll(strings.ToUpper(name), " ", "")
			compNameParsed := strings.ReplaceAll(strings.ToUpper(itemParsed.Name), " ", "")

			if compName == compNameParsed {
				foundInfo := false
				for infoIndex, infoParsed := range itemParsed.Infos {
					compInfo := strings.ReplaceAll(strings.ToUpper(info), " ", "")
					compInfoParsed := strings.ReplaceAll(strings.ToUpper(infoParsed.Name), " ", "")

					if compInfo == compInfoParsed {
						log.Trace().
							Str("compName", compName).
							Str("compNameParsed", compNameParsed).
							Str("compInfo", compInfo).
							Str("compInfoParsed", compInfoParsed).
							Msg("Same item and info names")
						dayParsed.Items[itemIndex].Infos[infoIndex].Start = append(infoParsed.Start, start)
						dayParsed.Items[itemIndex].Infos[infoIndex].End = append(infoParsed.End, end)

						foundInfo = true
						break
					}
				}

				if !foundInfo {
					log.Trace().
						Str("compName", compName).
						Str("compNameParsed", compNameParsed).
						Msg("Same item names, but no info")
					newInfo := Info{
						Name:  info,
						Start: make([]time.Time, 1),
						End:   make([]time.Time, 1),
					}
					newInfo.Start[0] = start
					newInfo.End[0] = end

					dayParsed.Items[itemIndex].Infos = append(itemParsed.Infos, newInfo)
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

func (day WeekdayParsed) Stringify() string {
	if len(day.Items) == 0 {
		return ""
	}

	spacer := "-------------------------"
	output := spacer + "\n\n**" + day.Name + ":**\n\n"

	for _, item := range day.Items {
		output += "--- **" + item.Name + "** ---\n"

		for _, info := range item.Infos {
			output += info.Name + "\n"

			for index := range info.Start {
				output += "**" + info.Start[index].Format("15:04") + "** - " + info.End[index].Format("15:04") + "\n"
			}
		}

		output += "\n"
	}
	output += spacer

	return output
}

func (week WeekParsed) Stringify() WeekOutput {
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
	return weekOutput
}

func generateAndParseWeek(items []*calendar.Event, namedDays config.NamedDays) WeekParsed {
	week := Week{
		Mon: Weekday{
			Name: namedDays.Monday,
		},
		Tue: Weekday{
			Name: namedDays.Tuesday,
		},
		Wed: Weekday{
			Name: namedDays.Wednesday,
		},
		Thu: Weekday{
			Name: namedDays.Thursday,
		},
		Fri: Weekday{
			Name: namedDays.Friday,
		},
		Sat: Weekday{
			Name: namedDays.Saturday,
		},
		Sun: Weekday{
			Name: namedDays.Sunday,
		},
	}

	week.Generate(items)
	weekParsed := week.Parse()
	log.Debug().
		Str("week", fmt.Sprintf("%v", weekParsed)).
		Msg("Parsed")

	return weekParsed
}

func updateCalendar(calendarId string, namedDays config.NamedDays, google config.Google) (WeekParsed, error) {
	week := WeekParsed{}

	log.Trace().
		Msg("Getting calendar via API")

	ctx := context.Background()
	calendarService, err := calendar.NewService(ctx, option.WithAPIKey(google.Token))
	if err != nil {
		return week, err
	}

	currentTime := time.Now().Format(time.RFC3339)
	weekTime := time.Now().AddDate(0, 0, 7).Format(time.RFC3339)

	if calendar, err := calendarService.Events.List(calendarId).ShowDeleted(false).
		SingleEvents(true).TimeMin(currentTime).TimeMax(weekTime).OrderBy("startTime").Do(); err != nil {
		return week, err
	} else {
		log.Debug().Msg("Decoded API response")
		week = generateAndParseWeek(calendar.Items, namedDays)
	}

	return week, nil
}

func outputDay(day string, webhookUrl string) error {
	if day != "" {
		if err := webhook.SendMessage(webhookUrl, webhook.Message{
			Content: day,
		}); err != nil {
			return err
		}
	}
	return nil
}

func outputWeek(channelId string, week WeekOutput) error {
	// messageObjects, err := ChannelMessages(channelId, 100, "", "", "")
	// if err != nil {
	// 	return err
	// }

	// messageIds := make([]string, 0)
	// for _, messageObject := range messageObjects {
	// 	messageIds = append(messageIds, messageObject.ID)
	// }

	// err = ChannelMessagesBulkDelete(channelId, messageIds)
	// if err != nil {
	// 	return err
	// }

	if err := outputDay(week.Mon, channelId); err != nil {
		return err
	}
	if err := outputDay(week.Tue, channelId); err != nil {
		return err
	}
	if err := outputDay(week.Wed, channelId); err != nil {
		return err
	}
	if err := outputDay(week.Thu, channelId); err != nil {
		return err
	}
	if err := outputDay(week.Fri, channelId); err != nil {
		return err
	}
	if err := outputDay(week.Sat, channelId); err != nil {
		return err
	}
	if err := outputDay(week.Sun, channelId); err != nil {
		return err
	}

	return nil
}

func getOldWeekOutput(webhookUrl string, namedDays config.NamedDays) (WeekOutput, error) {
	weekOutput := WeekOutput{}

	// messageObjects, err := discord.ChannelMessages(channelId, 100, "", "", "")
	// if err != nil {
	// 	return weekOutput, err
	// }

	// for _, messageObject := range messageObjects {
	// 	if strings.Contains(messageObject.Content, namedDays.Monday) {
	// 		weekOutput.Mon = messageObject.Content
	// 	}
	// 	if strings.Contains(messageObject.Content, namedDays.Tuesday) {
	// 		weekOutput.Tue = messageObject.Content
	// 	}
	// 	if strings.Contains(messageObject.Content, namedDays.Wednesday) {
	// 		weekOutput.Wed = messageObject.Content
	// 	}
	// 	if strings.Contains(messageObject.Content, namedDays.Thursday) {
	// 		weekOutput.Thu = messageObject.Content
	// 	}
	// 	if strings.Contains(messageObject.Content, namedDays.Friday) {
	// 		weekOutput.Fri = messageObject.Content
	// 	}
	// 	if strings.Contains(messageObject.Content, namedDays.Saturday) {
	// 		weekOutput.Sat = messageObject.Content
	// 	}
	// 	if strings.Contains(messageObject.Content, namedDays.Sunday) {
	// 		weekOutput.Sun = messageObject.Content
	// 	}
	// }

	return weekOutput, nil
}

func sameWeeks(weekA WeekOutput, weekB WeekOutput) bool {
	if weekA.Mon != weekB.Mon {
		return false
	}
	if weekA.Tue != weekB.Tue {
		return false
	}
	if weekA.Wed != weekB.Wed {
		return false
	}
	if weekA.Thu != weekB.Thu {
		return false
	}
	if weekA.Fri != weekB.Fri {
		return false
	}
	if weekA.Sat != weekB.Sat {
		return false
	}
	if weekA.Sun != weekB.Sun {
		return false
	}
	return true
}

func Update(ctx context.Context, conf *config.Config) {
	var worker conc.WaitGroup
	for _, calendarIterator := range conf.Calendars {
		calendarObject := calendarIterator
		worker.Go(func() {
			for {
				log.Debug().
					Str("name", calendarObject.Name).
					Msg("Updating calendar")

				week, err := updateCalendar(calendarObject.Id, conf.Days, conf.Google)
				if err != nil {
					log.Error().
						Err(err).
						Msg(fmt.Sprintf("Failed while updating calendar %s:", calendarObject.Name))
				}
				weekOutput := week.Stringify()
				weekOutputOld, err := getOldWeekOutput(calendarObject.Webhook, conf.Days)
				if err != nil {
					log.Error().
						Err(err).
						Msg(fmt.Sprintf("Failed while getting old calendar %s:", calendarObject.Name))
				}

				log.Debug().
					Str("new", fmt.Sprintf("%v", weekOutput)).
					Str("old", fmt.Sprintf("%v", weekOutputOld)).
					Msg("Comparing calendars")
				if sameWeeks(weekOutput, weekOutputOld) {
					log.Debug().
						Str("name", calendarObject.Name).
						Msg("Calendar is the same")
				} else {
					log.Trace().
						Str("name", calendarObject.Name).
						Msg("Outputting calendar")

					if err := outputWeek(calendarObject.Webhook, weekOutput); err != nil {
						log.Error().
							Err(err).
							Msg(fmt.Sprintf("Failed while outputting calendar %s:", calendarObject.Name))
					}
				}

				log.Trace().
					Str("name", calendarObject.Name).
					Msg("Sleeping calendar")

				if calendarObject.TimeBetweenChecks == 0 {
					time.Sleep(time.Hour * 3)
				} else {
					time.Sleep(time.Hour * time.Duration(calendarObject.TimeBetweenChecks))
				}
			}
		})
	}
	<-ctx.Done()
}
