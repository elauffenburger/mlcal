package calendar

import "github.com/gin-gonic/gin"

type CalendarRefresher interface {
	Refresh() error
}

func MakeRefreshCalendarEndpoint(refresher CalendarRefresher) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := refresher.Refresh(); err != nil {
			c.AbortWithError(500, err)
			return
		}

		c.String(200, "ok")
	}
}
