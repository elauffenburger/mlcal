package calendar

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/elauffenburger/musical-literacy-cal/pkg/mlcal"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func MakeGetCalendarHandler(calGetter Getter) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if e := recover(); e != nil {
				c.Error(fmt.Errorf("%v", e))
			}
		}()

		// Grab the calendar.
		cal, err := calGetter.Get()
		if err != nil || cal == nil {
			panic(errors.Wrapf(err, "error fetching calendar"))
		}

		// Grab the content type from the query params.
		var contentType string
		if c.DefaultQuery("text", "none") != "none" {
			contentType = "text"
		} else {
			contentType = "text/calendar"
		}

		// Check if the etag matches the request etag.
		etag := makeETag(cal, contentType)
		if reqETag := c.Request.Header.Get("If-None-Match"); reqETag != "" && reqETag == etag {
			c.Status(http.StatusNotModified)
			return
		}

		// Serialize the calendar and write the response.
		c.Header("Content-Type", contentType)
		c.Header("ETag", etag)
		c.String(200, "%s", cal.ToICS().Serialize())
	}
}

func makeETag(cal *mlcal.Calendar, contentType string) string {
	return fmt.Sprintf("%s_%d", strings.ReplaceAll(contentType, "/", "_"), cal.Date.UnixMilli())
}

type Getter interface {
	Get() (*mlcal.Calendar, error)
	GetICS() (string, error)
}

type Setter interface {
	Set(*mlcal.Calendar) error
}

type Cache interface {
	Getter
	Setter
}
