package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func respond(c *gin.Context, started time.Time, status int, payload gin.H) {
	if payload == nil {
		payload = gin.H{}
	}
	payload["time_ns"] = formatTimeNS(started)
	c.JSON(status, payload)
}

func fail(c *gin.Context, started time.Time, status int, message string) {
	respond(c, started, status, gin.H{"error": message})
}

func bindJSON(c *gin.Context, started time.Time, dst any) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		fail(c, started, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}
