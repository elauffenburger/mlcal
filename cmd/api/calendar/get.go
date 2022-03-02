package calendar

import (
	"fmt"
	"strings"

	"github.com/elauffenburger/musical-literacy-cal/cmd/api/resource"
	"github.com/elauffenburger/musical-literacy-cal/pkg/mlcal"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type calendarResource struct {
	etag        string
	contentType string
	cal         *mlcal.Calendar
}

func (r *calendarResource) ETag() string {
	return r.etag
}

func (r *calendarResource) Write(c *gin.Context) {
	c.Header("Content-Type", r.contentType)
	c.Header("ETag", r.etag)
	c.String(200, "%s", r.cal.ToICS().Serialize())
}

func MakeGetCalendarResource(calGetter Getter) resource.HandlerFunc {
	return func(c *gin.Context) (resource.Resource, error) {
		// Grab the calendar.
		cal, err := calGetter.Get()
		if err != nil || cal == nil {
			return nil, errors.Wrapf(err, "error fetching calendar")
		}

		// Grab the content type from the query params.
		var contentType string
		if c.DefaultQuery("text", "none") != "none" {
			contentType = "text"
		} else {
			contentType = "text/calendar"
		}

		// Derive an etag for the calendar.
		etag := fmt.Sprintf("%s_%d", strings.ReplaceAll(contentType, "/", "_"), cal.Date.UnixMilli())

		r := &calendarResource{
			cal:         cal,
			etag:        etag,
			contentType: contentType,
		}

		return r, nil
	}
}
