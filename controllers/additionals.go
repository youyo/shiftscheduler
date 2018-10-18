package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/youyo/shiftscheduler/models"
)

func GetAdditionals(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.GetAdditionals")

	status, message, err := models.GetAdditionals(c)
	c.JSON(status, response(message, err))
}

func PostAdditional(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.PostAdditional")

	var additional models.Additional
	if err := c.BindJSON(&additional); err != nil {
		c.JSON(400, response(nil, err))
		return
	}

	status, message, err := models.CreateAdditional(c, additional.RotationUuid, additional.Date, additional.Hour, additional.UserUuid)
	c.JSON(status, response(message, err))
}

func GetAdditional(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.GetAdditional")

	uuid := c.Param("uuid")

	status, message, err := models.GetAdditional(c, uuid)
	c.JSON(status, response(message, err))
}

func DeleteAdditional(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.DeleteAdditional")

	uuid := c.Param("uuid")

	status, message, err := models.DeleteAdditional(c, uuid)
	c.JSON(status, response(message, err))
}
