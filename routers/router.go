package routers

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/youyo/shiftscheduler/controllers"
	"github.com/youyo/shiftscheduler/middlewares"
)

var (
	Stage string = os.Getenv("STAGE")
)

func Setup() *gin.Engine {
	authMiddleware := middlewares.Jwt()
	r := gin.Default()

	r.GET("/", func(c *gin.Context) { c.Redirect(301, "/login") })
	r.GET("/login", func(c *gin.Context) {
		c.JSON(200, gin.H{"auth": "required"})
	})
	r.POST("/login", authMiddleware.LoginHandler)

	// required jwt
	api := r.Group("/api", authMiddleware.MiddlewareFunc())

	// refresh token
	api.GET("/refresh_token", authMiddleware.RefreshHandler)

	// users
	api.GET("/users", controllers.GetUsers)
	api.POST("/users", controllers.PostUser)
	api.GET("/users/:uuid", controllers.GetUser)
	api.DELETE("/users/:uuid", controllers.DeleteUser)
	//api.PATCH("/users/:uuid", controllers.PatchUser)

	// shifts
	api.GET("/shifts", controllers.GetShifts)
	api.POST("/shifts", controllers.PostShift)
	api.GET("/shifts/:uuid", controllers.GetShift)
	api.DELETE("/shifts/:uuid", controllers.DeleteShift)
	api.PATCH("/shifts/:uuid", controllers.PatchShift)

	// rotations
	api.GET("/rotations", controllers.GetRotations)
	api.POST("/rotations", controllers.PostRotation)
	api.GET("/rotations/:uuid", controllers.GetRotation)
	api.DELETE("/rotations/:uuid", controllers.DeleteRotation)
	//api.PATCH("/rotations/:uuid", controllers.PatchRotation)

	// rotation-details
	api.GET("/rotation-details", controllers.GetRotationDetails)
	api.POST("/rotation-details", controllers.PostRotationDetail)
	api.GET("/rotation-details/:uuid", controllers.GetRotationDetail)
	api.DELETE("/rotation-details/:uuid", controllers.DeleteRotationDetail)
	//api.PATCH("/rotation-details/:uuid", controllers.PatchRotationDetail)

	// overrides
	api.GET("/overrides", controllers.GetOverrides)
	api.POST("/overrides", controllers.PostOverride)
	api.GET("/overrides/:uuid", controllers.GetOverride)
	api.DELETE("/overrides/:uuid", controllers.DeleteOverride)
	//api.PATCH("/overrides/:uuid", controllers.PatchOverride)

	// additionals
	api.GET("/additionals", controllers.GetAdditionals)
	api.POST("/additionals", controllers.PostAdditional)
	api.GET("/additionals/:uuid", controllers.GetAdditional)
	api.DELETE("/additionals/:uuid", controllers.DeleteAdditional)
	//api.PATCH("/additionals/:uuid", controllers.PatchAdditional)

	// reduces
	api.GET("/reduces", controllers.GetReduces)
	api.POST("/reduces", controllers.PostReduce)
	api.GET("/reduces/:uuid", controllers.GetReduce)
	api.DELETE("/reduces/:uuid", controllers.DeleteReduce)
	//api.PATCH("/reduces/:uuid", controllers.PatchReduce)

	// schedules
	api.GET("/schedules/:rotation_uuid", func(c *gin.Context) {
		now := time.Now().In(time.Local)
		date := now.Format("2006-01-02")
		hour := now.Format("15")
		c.Redirect(301, "/api/schedules/"+c.Param("rotation_uuid")+"/"+date+"/"+hour)
	})
	api.GET("/schedules/:rotation_uuid/:date/:hour", controllers.GetSchedule)

	// calendars
	//api.GET("/calendars/:rotation_uuid/:date", controllers.GetCalendar)

	return r
}
