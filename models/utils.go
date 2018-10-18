package models

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"github.com/google/uuid"
	"github.com/hashicorp/logutils"
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

func GenUuid() string {
	return uuid.New().String()
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

func NewDbSession() (*dbr.Session, error) {
	dsn := BuildDsn(
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"))
	c, err := dbr.Open("mysql", dsn, nil)
	if err != nil {
		return nil, err
	}
	s := c.NewSession(nil)
	return s, nil
}

func BuildDsn(user, password, host, port, database string) string {
	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?parseTime=true&loc=Local"
	return dsn
}
