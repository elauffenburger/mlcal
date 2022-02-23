package mlcal

import (
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/google/uuid"
)

type Calendar struct {
	Date        time.Time
	Assignments []Assignment
}

type Assignment struct {
	Title string
	Due   time.Time
}

func (c *Calendar) ToICS() *ics.Calendar {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRefresh)
	cal.SetRefreshInterval("DURATION:P1D")

	now := time.Now()
	for _, a := range c.Assignments {
		e := cal.AddEvent(uuid.NewString())
		e.SetDtStampTime(a.Due)
		e.SetCreatedTime(now)
		e.SetAllDayStartAt(a.Due)
		e.SetSummary(a.Title)
		e.SetDescription(a.Title)
	}

	return cal
}
