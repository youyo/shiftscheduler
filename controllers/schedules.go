package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/youyo/shiftscheduler/models"
)

func GetSchedule(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.GetSchedule")

	rotationUuid := c.Param("rotation_uuid")
	date := c.Param("date")
	hour := c.Param("hour")

	status, message, err := models.GetSchedule(c, rotationUuid, date, hour)
	c.JSON(status, response(message, err))
}
