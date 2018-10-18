package models

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	OverridesTable string = "overrides"
)

type (
	Override struct {
		Uuid         string    `json:"uuid"`
		RotationUuid string    `json:"rotation_uuid"`
		Date         string    `json:"date"`
		Hour         string    `json:"hour"`
		UserUuid     string    `json:"user_uuid"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}

	Overrides []Override
)

func GetOverrides(c *gin.Context) (int, *Overrides, error) {
	LoggerDebug(c, "Called models.GetOverrides")

	sess, err := NewDbSession()
	if err != nil {
		return 500, nil, err
	}

	var overrides Overrides
	q := fmt.Sprintf("select * from %s", OverridesTable)
	if _, err := sess.SelectBySql(q).Load(&overrides); err != nil {
		LoggerError(c, fmt.Sprintf("failed to query. query: %s", q))
		return 400, nil, err
	}

	return 200, &overrides, nil
}

func CreateOverride(c *gin.Context, rotationUuid, date, hour, userUuid string) (int, string, error) {
	LoggerDebug(c, "Called models.CreateOverride")

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

	q := fmt.Sprintf("insert into %s (uuid,rotation_uuid,date,hour,user_uuid) values (?,?,?,?,?)", OverridesTable)
	if _, err := tx.InsertBySql(q, uuid, rotationUuid, date, hour, userUuid).Exec(); err != nil {
		tx.Rollback()
		LoggerError(c, fmt.Sprintf("failed to query. query: %s, uuid: %s, rotation_uuid: %s, date: %s, hour: %s, user_uuid: %s", q, uuid, rotationUuid, date, hour, userUuid))
		return 400, "", err
	}

	tx.Commit()

	return 201, "created", nil
}

func GetOverride(c *gin.Context, uuid string) (int, *Override, error) {
	LoggerDebug(c, "Called models.GetOverride")

	sess, err := NewDbSession()
	if err != nil {
		return 500, nil, err
	}

	var override Override
	q := fmt.Sprintf("select * from %s where uuid=?", OverridesTable)
	if _, err := sess.SelectBySql(q, uuid).Load(&override); err != nil {
		LoggerError(c, fmt.Sprintf("failed to query. query: %s, uuid: %s", q, uuid))
		return 400, nil, err
	}

	return 200, &override, nil
}

func DeleteOverride(c *gin.Context, uuid string) (int, string, error) {
	LoggerDebug(c, "Called models.DeleteOverride")

	sess, err := NewDbSession()
	if err != nil {
		return 500, "", err
	}

	tx, err := sess.Begin()
	if err != nil {
		return 500, "", err
	}
	defer tx.RollbackUnlessCommitted()

	q := fmt.Sprintf("delete from %s where uuid=?", OverridesTable)
	if _, err := tx.DeleteBySql(q, uuid).Exec(); err != nil {
		tx.Rollback()
		LoggerError(c, fmt.Sprintf("failed to query. query: %s, uuid: %s", q, uuid))
		return 400, "", err
	}

	tx.Commit()

	return 204, "deleted", nil
}
