package models

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ReducesTable string = "reduces"
)

type (
	Reduce struct {
		Uuid         string    `json:"uuid"`
		RotationUuid string    `json:"rotation_uuid"`
		Date         string    `json:"date"`
		Hour         string    `json:"hour"`
		UserUuid     string    `json:"user_uuid"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}

	Reduces []Reduce
)

func GetReduces(c *gin.Context) (int, *Reduces, error) {
	LoggerDebug(c, "Called models.GetReduces")

	sess, err := NewDbSession()
	if err != nil {
		return 500, nil, err
	}

	var reduces Reduces
	q := fmt.Sprintf("select * from %s", ReducesTable)
	if _, err := sess.SelectBySql(q).Load(&reduces); err != nil {
		LoggerError(c, fmt.Sprintf("failed to query. query: %s", q))
		return 400, nil, err
	}

	return 200, &reduces, nil
}

func CreateReduce(c *gin.Context, rotationUuid, date, hour, userUuid string) (int, string, error) {
	LoggerDebug(c, "Called models.CreateReduce")

	sess, err := NewDbSession()
	if err != nil {
		return 500, "", err
	}

	tx, err := sess.Begin()
	if err != nil {
		return 500, "", err
	}
	defer tx.RollbackUnlessCommitted()

	uuid := GenUuid()

	q := fmt.Sprintf("insert into %s (uuid,rotation_uuid,date,hour,user_uuid) values (?,?,?,?,?)", ReducesTable)
	if _, err := tx.InsertBySql(q, uuid, rotationUuid, date, hour, userUuid).Exec(); err != nil {
		tx.Rollback()
		LoggerError(c, fmt.Sprintf("failed to query. query: %s, uuid: %s, rotation_uuid: %s, date: %s, hour: %s, user_uuid: %s", q, uuid, rotationUuid, date, hour, userUuid))
		return 400, "", err
	}

	tx.Commit()

	return 201, "created", nil
}

func GetReduce(c *gin.Context, uuid string) (int, *Reduce, error) {
	LoggerDebug(c, "Called models.GetReduce")

	sess, err := NewDbSession()
	if err != nil {
		return 500, nil, err
	}

	var reduce Reduce
	q := fmt.Sprintf("select * from %s where uuid=?", ReducesTable)
	if _, err := sess.SelectBySql(q, uuid).Load(&reduce); err != nil {
		LoggerError(c, fmt.Sprintf("failed to query. query: %s, uuid: %s", q, uuid))
		return 400, nil, err
	}

	return 200, &reduce, nil
}

func DeleteReduce(c *gin.Context, uuid string) (int, string, error) {
	LoggerDebug(c, "Called models.DeleteReduce")

	sess, err := NewDbSession()
	if err != nil {
		return 500, "", err
	}

	tx, err := sess.Begin()
	if err != nil {
		return 500, "", err
	}
	defer tx.RollbackUnlessCommitted()

	q := fmt.Sprintf("delete from %s where uuid=?", ReducesTable)
	if _, err := tx.DeleteBySql(q, uuid).Exec(); err != nil {
		tx.Rollback()
		LoggerError(c, fmt.Sprintf("failed to query. query: %s, uuid: %s", q, uuid))
		return 400, "", err
	}

	tx.Commit()

	return 204, "deleted", nil
}
