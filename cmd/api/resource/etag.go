package resource

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Resource interface {
	ETag() string
	Write(*gin.Context)
}

type HandlerFunc func(*gin.Context) (Resource, error)

func Handler(getResource HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := getResource(c)
		if err != nil {
			c.Error(err)
			return
		}

		reqTag := c.Request.Header.Get("If-None-Match")
		if reqTag != "" && reqTag == res.ETag() {
			c.Status(http.StatusNotModified)
			return
		}

		res.Write(c)

		c.Next()
	}
}
