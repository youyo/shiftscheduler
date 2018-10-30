package models

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	UsersTable string = "users"
)

type (
	User struct {
		Uuid        string    `json:"uuid"`
		Name        string    `json:"name" binding:"required"`
		PhoneNumber string    `json:"phone_number" binding:"required"`
		CreatedAt   time.Time `json:"created_at,omitempty"`
		UpdatedAt   time.Time `json:"updated_at,omitempty"`
	}

	Users []User
)

func GetUsers(c *gin.Context) (int, *Users, error) {
	LoggerDebug(c, "Called models.GetUsers")

	sess, err := NewDbSession()
	if err != nil {
		return 500, nil, err
	}

	var users Users
	q := fmt.Sprintf("select * from %s", UsersTable)
	if _, err := sess.SelectBySql(q).Load(&users); err != nil {
		LoggerError(c, fmt.Sprintf("failed to query. query: %s", q))
		return 400, nil, err
	}

	return 200, &users, nil
}

func CreateUser(c *gin.Context, name, phoneNumber string) (int, string, error) {
	LoggerDebug(c, "Called models.CreateUser")

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

	q := fmt.Sprintf("insert into %s (uuid,name,phone_number) values (?,?,?)", UsersTable)
	if _, err := tx.InsertBySql(q, uuid, name, phoneNumber).Exec(); err != nil {
		tx.Rollback()
		LoggerError(c, fmt.Sprintf("failed to query. query: %s, uuid: %s, name: %s, phone_number: %s", q, uuid, name, phoneNumber))
		return 400, "", err
	}

	tx.Commit()

	return 201, "created", nil
}

func GetUser(c *gin.Context, uuid string) (int, *User, error) {
	LoggerDebug(c, "Called models.GetUser")

	sess, err := NewDbSession()
	if err != nil {
		return 500, nil, err
	}

	var user User
	q := fmt.Sprintf("select * from %s where uuid=?", UsersTable)
	if _, err := sess.SelectBySql(q, uuid).Load(&user); err != nil {
		LoggerError(c, fmt.Sprintf("failed to query. query: %s, uuid: %s", q, uuid))
		return 400, nil, err
	}

	return 200, &user, nil
}

func DeleteUser(c *gin.Context, uuid string) (int, string, error) {
	LoggerDebug(c, "Called models.DeleteUser")

	sess, err := NewDbSession()
	if err != nil {
		return 500, "", err
	}

	tx, err := sess.Begin()
	if err != nil {
		return 500, "", err
	}
	defer tx.RollbackUnlessCommitted()

	q := fmt.Sprintf("delete from %s where uuid=?", UsersTable)
	if _, err := tx.DeleteBySql(q, uuid).Exec(); err != nil {
		tx.Rollback()
		LoggerError(c, fmt.Sprintf("failed to query. query: %s, uuid: %s", q, uuid))
		return 400, "", err
	}

	tx.Commit()

	return 204, "deleted", nil
}
