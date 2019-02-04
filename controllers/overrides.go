package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/youyo/shiftscheduler/models"
)

func GetOverrides(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.GetOverrides")

	allRecords := c.DefaultQuery("all", "false")
	status, message, err := models.GetOverrides(c, allRecords)
	c.JSON(status, response(message, err))
}

func PostOverride(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.PostOverride")

	var override models.Override
	if err := c.BindJSON(&override); err != nil {
		c.JSON(400, response(nil, err))
		return
	}

	status, message, err := models.CreateOverride(c, override.RotationUuid, override.Date, override.Hour, override.UserUuid)
	c.JSON(status, response(message, err))
}

func GetOverride(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.GetOverride")

	uuid := c.Param("uuid")

	status, message, err := models.GetOverride(c, uuid)
	c.JSON(status, response(message, err))
}

func DeleteOverride(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.DeleteOverride")

	uuid := c.Param("uuid")

	status, message, err := models.DeleteOverride(c, uuid)
	c.JSON(status, response(message, err))
}
