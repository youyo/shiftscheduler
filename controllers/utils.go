package controllers

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hashicorp/logutils"
)

type (
	Response struct {
		Response interface{} `json:"response"`
		Error    interface{} `json:"error"`
	}
)

func init() {
	filter := &logutils.LevelFilter{
		Levels: []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		Writer: os.Stderr,
	}

	if os.Getenv("GIN_MODE") == "release" {
		filter.MinLevel = logutils.LogLevel("INFO")
	} else {
		filter.MinLevel = logutils.LogLevel("DEBUG")
	}

	//log.SetPrefix("[APP] ")
	log.SetOutput(filter)
}

func SetRequestId(c *gin.Context) {
	c.Set("request_id", uuid.New().String())
}

func GetRequestId(c *gin.Context) string {
	return c.GetString("request_id")
}

func printMessage(loglevel, requestId, message string) string {
	return fmt.Sprintf("[%s] request_id: %s, message: %s", loglevel, requestId, message)
}

func LoggerDebug(c *gin.Context, message string) {
	log.Printf(printMessage("DEBUG", GetRequestId(c), message))
}

func LoggerInfo(c *gin.Context, message string) {
	log.Printf(printMessage("INFO", GetRequestId(c), message))
}

func LoggerWarn(c *gin.Context, message string) {
	log.Printf(printMessage("WARN", GetRequestId(c), message))
}

func LoggerError(c *gin.Context, message string) {
	log.Printf(printMessage("ERROR", GetRequestId(c), message))
}

func response(res interface{}, err error) Response {
	r := Response{
		Response: res,
		Error:    err,
	}
	if err != nil {
		r.Error = fmt.Sprintf("%v", err)
	}
	return r
}
