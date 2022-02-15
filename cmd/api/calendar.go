package main

import (
	"fmt"
	"os"

	"github.com/elauffenburger/musical-literacy-cal/pkg/mlcal"
	"github.com/gin-gonic/gin"
)

func makeGetCalendarHandler(icsCal string) gin.HandlerFunc {
	return func(c *gin.Context) {
		getCalendar(c, icsCal)
	}
}

func getCalendar(c *gin.Context, icsCal string) {
	c.String(200, "%s", icsCal)
}

func loadICSCal(email, password, classID string) (string, error) {
	client, err := mlcal.NewClient(email, password, classID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cal, err := client.Get()
	if err != nil {
		return "", err
	}

	return cal.ToICS().Serialize(), nil
}
