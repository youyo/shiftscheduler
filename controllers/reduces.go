package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/youyo/shiftscheduler/models"
)

func GetReduces(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.GetReduces")

	status, message, err := models.GetReduces(c)
	c.JSON(status, response(message, err))
}

func PostReduce(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.PostReduce")

	var reduce models.Reduce
	if err := c.BindJSON(&reduce); err != nil {
		c.JSON(400, response(nil, err))
		return
	}

	status, message, err := models.CreateReduce(c, reduce.RotationUuid, reduce.Date, reduce.Hour, reduce.UserUuid)
	c.JSON(status, response(message, err))
}

func GetReduce(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.GetReduce")

	uuid := c.Param("uuid")

	status, message, err := models.GetReduce(c, uuid)
	c.JSON(status, response(message, err))
}

func DeleteReduce(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.DeleteReduce")

	uuid := c.Param("uuid")

	status, message, err := models.DeleteReduce(c, uuid)
	c.JSON(status, response(message, err))
}
