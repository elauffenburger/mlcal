package main

import (
	"github.com/elauffenburger/musical-literacy-cal/cmd/api/calendar"
	"github.com/elauffenburger/musical-literacy-cal/cmd/api/resource"
	"github.com/gin-gonic/gin"
)

func addCalendarEndpoints(srv *gin.Engine, calRefresher interface {
	calendar.Getter
	calendar.CalendarRefresher
}) error {
	srv.GET("/calendar", resource.Handler(calendar.MakeGetCalendarResource(calRefresher)))
	srv.GET("/calendar/refresh", calendar.MakeRefreshCalendarEndpoint(calRefresher))

	return nil
}
