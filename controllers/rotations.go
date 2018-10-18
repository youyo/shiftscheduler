package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/youyo/shiftscheduler/models"
)

func GetRotations(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.GetRotations")

	status, message, err := models.GetRotations(c)
	c.JSON(status, response(message, err))
}

func PostRotation(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.PostRotation")

	var rotation models.Rotation
	if err := c.BindJSON(&rotation); err != nil {
		c.JSON(400, response(nil, err))
		return
	}

	status, message, err := models.CreateRotation(c, rotation.Name, rotation.StartDate)
	c.JSON(status, response(message, err))
}

func GetRotation(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.GetRotation")

	uuid := c.Param("uuid")

	status, message, err := models.GetRotation(c, uuid)
	c.JSON(status, response(message, err))
}

func DeleteRotation(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.DeleteRotation")

	uuid := c.Param("uuid")

	status, message, err := models.DeleteRotation(c, uuid)
	c.JSON(status, response(message, err))
}
