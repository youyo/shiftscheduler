package models

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	RotationsTable string = "rotations"
)

type (
	Rotation struct {
		Uuid      string    `json:"uuid"`
		Name      string    `json:"name"`
		StartDate string    `json:"start_date"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	Rotations []Rotation
)

func GetRotations(c *gin.Context) (int, *Rotations, error) {
	LoggerDebug(c, "Called models.GetRotations")

	sess, err := NewDbSession()
	if err != nil {
		return 500, nil, err
	}

	var rotations Rotations
	q := fmt.Sprintf("select * from %s", RotationsTable)
	if _, err := sess.SelectBySql(q).Load(&rotations); err != nil {
		LoggerError(c, fmt.Sprintf("failed to query. query: %s", q))
		return 400, nil, err
	}

	return 200, &rotations, nil
}

func CreateRotation(c *gin.Context, name, StartDate string) (int, string, error) {
	LoggerDebug(c, "Called models.CreateRotation")

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

	q := fmt.Sprintf("insert into %s (uuid,name,start_date) values (?,?,?)", RotationsTable)
	if _, err := tx.InsertBySql(q, uuid, name, StartDate).Exec(); err != nil {
		tx.Rollback()
		LoggerError(c, fmt.Sprintf("failed to query. query: %s, uuid: %s, name: %s, start_date: %s", q, uuid, name, StartDate))
		return 400, "", err
	}

	tx.Commit()

	return 201, "created", nil
}

func GetRotation(c *gin.Context, uuid string) (int, *Rotation, error) {
	LoggerDebug(c, "Called models.GetRotation")

	sess, err := NewDbSession()
	if err != nil {
		return 500, nil, err
	}

	var rotation Rotation
	q := fmt.Sprintf("select * from %s where uuid=?", RotationsTable)
	if _, err := sess.SelectBySql(q, uuid).Load(&rotation); err != nil {
		LoggerError(c, fmt.Sprintf("failed to query. query: %s, uuid: %s", q, uuid))
		return 400, nil, err
	}

	return 200, &rotation, nil
}

func DeleteRotation(c *gin.Context, uuid string) (int, string, error) {
	LoggerDebug(c, "Called models.DeleteRotation")

	sess, err := NewDbSession()
	if err != nil {
		return 500, "", err
	}

	tx, err := sess.Begin()
	if err != nil {
		return 500, "", err
	}
	defer tx.RollbackUnlessCommitted()

	q := fmt.Sprintf("delete from %s where uuid=?", RotationsTable)
	if _, err := tx.DeleteBySql(q, uuid).Exec(); err != nil {
		tx.Rollback()
		LoggerError(c, fmt.Sprintf("failed to query. query: %s, uuid: %s", q, uuid))
		return 400, "", err
	}

	tx.Commit()

	return 204, "deleted", nil
}
