package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocraft/dbr"
)

const (
	DateTimeFormat string = "2006-01-02-15-04"
	DateFormat     string = "2006-01-02"
	HourFormat     string = "15"
)

func GetSchedule(c *gin.Context, rotationUuid, date, hour string) (int, Users, error) {
	LoggerDebug(c, "Called models.GetSchedule")

	// create db session
	LoggerDebug(c, "create db session")
	sess, err := NewDbSession()
	if err != nil {
		return 500, nil, err
	}

	status, users, err := QuerySchedule(c, sess, rotationUuid, date, hour)
	return status, users, err
}

func QuerySchedule(c *gin.Context, sess *dbr.Session, rotationUuid, date, hour string) (int, Users, error) {
	LoggerDebug(c, "Called models.QuerySchedule")

	// check overrides
	status, overrideUserUuids, err := QueryOverride(c, sess, rotationUuid, date, hour)
	if err != nil {
		LoggerError(c, fmt.Sprintf("check override error"))
		return status, nil, err
	}

	if len(overrideUserUuids) != 0 {
		var users Users
		q := fmt.Sprintf("select uuid,name,phone_number from %s where uuid in ?", UsersTable)
		if _, err := sess.SelectBySql(q, overrideUserUuids).Load(&users); err != nil {
			LoggerError(c, fmt.Sprintf("failed to query. query: %s, uuid: %v, err: %v", q, overrideUserUuids, err))
			return 400, nil, err
		}
		return 200, users, nil
	}

	// get rotation start_date
	LoggerDebug(c, "get rotation start_date")
	var startDate string
	q := fmt.Sprintf("select start_date from %s where uuid=?", RotationsTable)
	if _, err := sess.SelectBySql(q, rotationUuid).Load(&startDate); err != nil {
		LoggerError(c, fmt.Sprintf("failed to query. query: %s, uuid: %s", q, rotationUuid))
		return 400, nil, err
	}

	// スケジュール開始日と取得したい日の差分
	LoggerDebug(c, "get duration")
	dateT, err := time.ParseInLocation(DateFormat, date, time.Local)
	if err != nil {
		return 500, nil, err
	}

	startDateT, err := time.ParseInLocation(DateFormat, startDate, time.Local)
	if err != nil {
		return 500, nil, err
	}

	duration := dateT.Sub(startDateT)
	if duration < 0 {
		var users Users
		return 200, users, nil
	}

	// duration の総時間から日数を取得
	LoggerDebug(c, "get days")
	days := int(duration.Hours()) / 24

	// 最大 order_id 取得
	// 週数
	LoggerDebug(c, "get max order_id")
	var maxOrderId int
	q = fmt.Sprintf("select max(order_id) as order_id from %s where rotation_uuid=? group by rotation_uuid", RotationDetailsTable)
	if _, err := sess.SelectBySql(q, rotationUuid).Load(&maxOrderId); err != nil {
		LoggerError(c, fmt.Sprintf("failed to query. query: %s, rotation_uuid: %s", q, rotationUuid))
		return 400, nil, err
	}

	// シフトの全日数を取得
	LoggerDebug(c, "get total days")
	totalDays := maxOrderId * 7

	// 曜日取得
	LoggerDebug(c, "get week")
	var week string
	switch days % 7 {
	case 0:
		week = "monday_"
	case 1:
		week = "tuesday_"
	case 2:
		week = "wednesday_"
	case 3:
		week = "thursday_"
	case 4:
		week = "friday_"
	case 5:
		week = "saturday_"
	case 6:
		week = "sunday_"
	default:
		return 500, nil, errors.New("Unmatch days")
	}

	// order_id 確定
	LoggerDebug(c, "get order_id")
	for days > totalDays {
		days = days - totalDays
	}
	orderId := 0
	for orderIdStart := 1; orderIdStart <= maxOrderId; orderIdStart++ {
		if days <= orderIdStart*7 {
			orderId = orderIdStart
			break
		}
	}

	// shift_uuids 取得
	LoggerDebug(c, "get shift_uuids")
	var shiftUuids []string
	q = fmt.Sprintf("select shift_uuid from %s where rotation_uuid=? and order_id=?", RotationDetailsTable)
	if _, err := sess.SelectBySql(q, rotationUuid, orderId).Load(&shiftUuids); err != nil {
		LoggerError(c, fmt.Sprintf("failed to query. query: %s, rotation_uuid: %s, order_id: %d", q, rotationUuid, orderId))
		return 400, nil, err
	}

	// 有効な user_uuid 取得
	LoggerDebug(c, "get user uuids")
	var userUuids []string
	q = fmt.Sprintf("select user_uuid from %s where %s=1 and uuid in ?", ShiftsTable, week+hour+"00")
	if _, err := sess.SelectBySql(q, shiftUuids).Load(&userUuids); err != nil {
		LoggerError(c, fmt.Sprintf("failed to query. query: %s, uuid: %v", q, shiftUuids))
		return 400, nil, err
	}

	if len(userUuids) == 0 {
		err := errors.New("matched user is not exist")
		LoggerWarn(c, fmt.Sprintf("%v", err))
		return 400, nil, err
	}

	// ユーザー情報取得
	LoggerDebug(c, "get user info")
	var users Users
	q = fmt.Sprintf("select uuid,name,phone_number from %s where uuid in ?", UsersTable)
	if _, err := sess.SelectBySql(q, userUuids).Load(&users); err != nil {
		LoggerError(c, fmt.Sprintf("failed to query. query: %s, uuid: %v", q, userUuids))
		return 400, nil, err
	}

	return 200, users, nil
}

func QueryOverride(c *gin.Context, sess *dbr.Session, rotationUuid, date, hour string) (int, []string, error) {
	LoggerDebug(c, "Called models.QueryOverride")

	var userUuids []string
	q := fmt.Sprintf("select user_uuid from %s where rotation_uuid=? and date=? and hour=?", OverridesTable)
	if _, err := sess.SelectBySql(q, rotationUuid, date, hour).Load(&userUuids); err != nil {
		LoggerError(c, fmt.Sprintf("failed to query. query: %s, rotation_uuid: %s, date: %s, hour: %s", q, rotationUuid, date, hour))
		return 400, nil, err
	}

	return 200, userUuids, nil
}
