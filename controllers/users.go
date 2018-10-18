package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/youyo/shiftscheduler/models"
)

func GetUsers(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.GetUsers")

	status, message, err := models.GetUsers(c)
	c.JSON(status, response(message, err))
}

func PostUser(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.PostUser")

	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, response(nil, err))
		return
	}

	status, message, err := models.CreateUser(c, user.Name, user.PhoneNumber)
	c.JSON(status, response(message, err))
}

func GetUser(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.GetUser")

	uuid := c.Param("uuid")

	status, message, err := models.GetUser(c, uuid)
	c.JSON(status, response(message, err))
}

func DeleteUser(c *gin.Context) {
	SetRequestId(c)
	LoggerDebug(c, "Called controllers.DeleteUser")

	uuid := c.Param("uuid")

	status, message, err := models.DeleteUser(c, uuid)
	c.JSON(status, response(message, err))
}
