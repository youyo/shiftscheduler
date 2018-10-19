package models

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type (
	CalendarEvent struct {
		Title string `json:"title"`
		Start string `json:"start"`
	}
	CalendarEvents []CalendarEvent
)

func GetCalendar(c *gin.Context, rotationUuid, date string) (int, CalendarEvents, error) {
	LoggerDebug(c, "Called models.GetCalendar")
	var calendarEvents CalendarEvents

	// first day
	firstDay := 1

	// get last day
	lastDay, err := GetLastDay(date)
	if err != nil {
		return 400, nil, err
	}

	// get all days
	days := GetAllDays(date, firstDay, lastDay)

	// create db session
	sess, err := NewDbSession()
	if err != nil {
		return 500, nil, err
	}

	// get schedules
	for _, day := range days {
		hours := []string{
			"00", "01", "02", "03", "04", "05", "06", "07", "08", "09",
			"10", "11", "12", "13", "14", "15", "16", "17", "18", "19",
			"20", "21", "22", "23",
		}
		for _, hour := range hours {
			status, users, err := QuerySchedule(c, sess, rotationUuid, day, hour)
			if err != nil {
				return status, nil, err
			}
			if len(users) == 0 {
				continue
			}

			startDate, err := time.Parse("2006-01-02-15", day+"-"+hour)
			if err != nil {
				return 500, nil, err
			}

			for _, user := range users {
				cal := CalendarEvent{
					Title: user.Name,
					Start: startDate.Format("2006-01-02T15:04:05"),
				}
				calendarEvents = append(calendarEvents, cal)
			}
		}
	}

	return 200, calendarEvents, nil
}

func GetLastDay(date string) (int, error) {
	t, err := time.ParseInLocation("2006-01", date, time.Local)
	if err != nil {
		return 0, err
	}

	lastDayTime := t.AddDate(0, 1, 0).AddDate(0, 0, -1)
	lastDay, err := strconv.Atoi(lastDayTime.Format("02"))
	return lastDay, err
}

func GetAllDays(date string, firstDay, lastDay int) []string {
	var days []string
	for day := firstDay; day < 10; day++ {
		s := date + "-0" + strconv.Itoa(day)
		days = append(days, s)
	}
	for day := 10; day <= lastDay; day++ {
		s := date + "-" + strconv.Itoa(day)
		days = append(days, s)
	}
	return days
}
