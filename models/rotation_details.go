package models

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	RotationDetailsTable string = "rotation_details"
)

type (
	RotationDetail struct {
		Uuid         string    `json:"uuid"`
		RotationUuid string    `json:"rotation_uuid"`
		OrderId      int       `json:"order_id"`
		ShiftUuid    string    `json:"shift_uuid"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}

	RotationDetails []RotationDetail
)

func GetRotationDetails(c *gin.Context) (int, *RotationDetails, error) {
	LoggerDebug(c, "Called models.GetRotationDetails")

	sess, err := NewDbSession()
	if err != nil {
		return 500, nil, err
	}

	var rotationDetails RotationDetails
	q := fmt.Sprintf("select * from %s", RotationDetailsTable)
	if _, err := sess.SelectBySql(q).Load(&rotationDetails); err != nil {
		LoggerError(c, fmt.Sprintf("failed to query. query: %s", q))
		return 400, nil, err
	}

	return 200, &rotationDetails, nil
}

func CreateRotationDetail(c *gin.Context, rotationUuid, shiftUuid string, orderId int) (int, string, error) {
	LoggerDebug(c, "Called models.CreateRotationDetail")

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

	q := fmt.Sprintf("insert into %s (uuid,rotation_uuid,order_id,shift_uuid) values (?,?,?)", RotationDetailsTable)
	if _, err := tx.InsertBySql(q, uuid, rotationUuid, orderId, shiftUuid).Exec(); err != nil {
		tx.Rollback()
		LoggerError(c, fmt.Sprintf("failed to query. query: %s, uuid: %s, rotation_uuid: %s, order_id: %d, shift_uuid: %s", q, uuid, rotationUuid, orderId, shiftUuid))
		return 400, "", err
	}

	tx.Commit()

	return 201, "created", nil
}

func GetRotationDetail(c *gin.Context, uuid string) (int, *RotationDetail, error) {
	LoggerDebug(c, "Called models.GetRotationDetail")

	sess, err := NewDbSession()
	if err != nil {
		return 500, nil, err
	}

	var rotationDetail RotationDetail
	q := fmt.Sprintf("select * from %s where uuid=?", RotationDetailsTable)
	if _, err := sess.SelectBySql(q, uuid).Load(&rotationDetail); err != nil {
		LoggerError(c, fmt.Sprintf("failed to query. query: %s, uuid: %s", q, uuid))
		return 400, nil, err
	}

	return 200, &rotationDetail, nil
}

func DeleteRotationDetail(c *gin.Context, uuid string) (int, string, error) {
	LoggerDebug(c, "Called models.DeleteRotationDetail")

	sess, err := NewDbSession()
	if err != nil {
		return 500, "", err
	}

	tx, err := sess.Begin()
	if err != nil {
		return 500, "", err
	}
	defer tx.RollbackUnlessCommitted()

	q := fmt.Sprintf("delete from %s where uuid=?", RotationDetailsTable)
	if _, err := tx.DeleteBySql(q, uuid).Exec(); err != nil {
		tx.Rollback()
		LoggerError(c, fmt.Sprintf("failed to query. query: %s, uuid: %s", q, uuid))
		return 400, "", err
	}

	tx.Commit()

	return 204, "deleted", nil
}
