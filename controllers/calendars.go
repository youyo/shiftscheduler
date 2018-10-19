package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/youyo/shiftscheduler/models"
)

func GetCalendar(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.GetCalendar")

	rotationUuid := c.Param("rotation_uuid")
	date := c.Param("date")

	status, message, err := models.GetCalendar(c, rotationUuid, date)
	c.JSON(status, response(message, err))
}
