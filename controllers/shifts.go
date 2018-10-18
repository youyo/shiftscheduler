package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/youyo/shiftscheduler/models"
)

func GetShifts(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.GetShifts")

	status, message, err := models.GetShifts(c)
	c.JSON(status, response(message, err))
}

func PostShift(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.PostShift")

	var shift models.Shift
	if err := c.BindJSON(&shift); err != nil {
		c.JSON(400, response(nil, err))
		return
	}

	status, message, err := models.CreateShift(c, shift.Name, shift.UserUuid)
	c.JSON(status, response(message, err))
}

func GetShift(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.GetShift")

	uuid := c.Param("uuid")

	status, message, err := models.GetShift(c, uuid)
	c.JSON(status, response(message, err))
}

func DeleteShift(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.DeleteShift")

	uuid := c.Param("uuid")

	status, message, err := models.DeleteShift(c, uuid)
	c.JSON(status, response(message, err))
}

func PatchShift(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.PatchShift")

	var shiftDetail models.ShiftDetail
	if err := c.BindJSON(&shiftDetail); err != nil {
		c.JSON(400, response(nil, err))
		return
	}

	uuid := c.Param("uuid")

	status, message, err := models.PatchShift(c, uuid, shiftDetail)
	c.JSON(status, response(message, err))
}
