package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/youyo/shiftscheduler/models"
)

func GetRotationDetails(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.GetRotationDetails")

	status, message, err := models.GetRotationDetails(c)
	c.JSON(status, response(message, err))
}

func PostRotationDetail(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.PostRotationDetail")

	var rotationDetail models.RotationDetail
	if err := c.BindJSON(&rotationDetail); err != nil {
		c.JSON(400, response(nil, err))
		return
	}

	status, message, err := models.CreateRotationDetail(c, rotationDetail.RotationUuid, rotationDetail.ShiftUuid, rotationDetail.OrderId)
	c.JSON(status, response(message, err))
}

func GetRotationDetail(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.GetRotationDetail")

	uuid := c.Param("uuid")

	status, message, err := models.GetRotationDetail(c, uuid)
	c.JSON(status, response(message, err))
}

func DeleteRotationDetail(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.DeleteRotationDetail")

	uuid := c.Param("uuid")

	status, message, err := models.DeleteRotationDetail(c, uuid)
	c.JSON(status, response(message, err))
}
