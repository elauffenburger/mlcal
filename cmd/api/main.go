package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load the cal on startup.
	icsCal, err := loadICSCal(
		os.Getenv("ML_EMAIL"),
		os.Getenv("ML_PASSWORD"),
		os.Getenv("ML_COURSE"),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Start server.
	r := gin.Default()
	r.GET("/calendar", makeGetCalendarHandler(icsCal))

	r.Run()
}
