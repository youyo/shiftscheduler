package models

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	AdditionalsTable string = "additionals"
)

type (
	Additional struct {
		Uuid         string    `json:"uuid"`
		RotationUuid string    `json:"rotation_uuid" binding:"required"`
		Date         string    `json:"date" binding:"required"`
		Hour         string    `json:"hour" binding:"required"`
		UserUuid     string    `json:"user_uuid" binding:"required"`
		CreatedAt    time.Time `json:"created_at,omitempty"`
		UpdatedAt    time.Time `json:"updated_at,omitempty"`
	}

	Additionals []Additional
)

func GetAdditionals(c *gin.Context) (int, *Additionals, error) {
	LoggerDebug(c, "Called models.GetAdditionals")

	sess, err := NewDbSession()
	if err != nil {
		return 500, nil, err
	}

	var additionals Additionals
	q := fmt.Sprintf("select * from %s", AdditionalsTable)
	if _, err := sess.SelectBySql(q).Load(&additionals); err != nil {
		LoggerError(c, fmt.Sprintf("failed to query. query: %s", q))
		return 400, nil, err
	}

	return 200, &additionals, nil
}

func CreateAdditional(c *gin.Context, rotationUuid, date, hour, userUuid string) (int, string, error) {
	LoggerDebug(c, "Called models.CreateAdditional")

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

	q := fmt.Sprintf("insert into %s (uuid,rotation_uuid,date,hour,user_uuid) values (?,?,?,?,?)", AdditionalsTable)
	if _, err := tx.InsertBySql(q, uuid, rotationUuid, date, hour, userUuid).Exec(); err != nil {
		tx.Rollback()
		LoggerError(c, fmt.Sprintf("failed to query. query: %s, uuid: %s, rotation_uuid: %s, date: %s, hour: %s, user_uuid: %s", q, uuid, rotationUuid, date, hour, userUuid))
		return 400, "", err
	}

	tx.Commit()

	return 201, "created", nil
}

func GetAdditional(c *gin.Context, uuid string) (int, *Additional, error) {
	LoggerDebug(c, "Called models.GetAdditional")

	sess, err := NewDbSession()
	if err != nil {
		return 500, nil, err
	}

	var additional Additional
	q := fmt.Sprintf("select * from %s where uuid=?", AdditionalsTable)
	if _, err := sess.SelectBySql(q, uuid).Load(&additional); err != nil {
		LoggerError(c, fmt.Sprintf("failed to query. query: %s, uuid: %s", q, uuid))
		return 400, nil, err
	}

	return 200, &additional, nil
}

func DeleteAdditional(c *gin.Context, uuid string) (int, string, error) {
	LoggerDebug(c, "Called models.DeleteAdditional")

	sess, err := NewDbSession()
	if err != nil {
		return 500, "", err
	}

	tx, err := sess.Begin()
	if err != nil {
		return 500, "", err
	}
	defer tx.RollbackUnlessCommitted()

	q := fmt.Sprintf("delete from %s where uuid=?", AdditionalsTable)
	if _, err := tx.DeleteBySql(q, uuid).Exec(); err != nil {
		tx.Rollback()
		LoggerError(c, fmt.Sprintf("failed to query. query: %s, uuid: %s", q, uuid))
		return 400, "", err
	}

	tx.Commit()

	return 204, "deleted", nil
}
