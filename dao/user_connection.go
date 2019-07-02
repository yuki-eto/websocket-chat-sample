package dao

import (
	"fmt"

	"github.com/cornelk/hashmap"
	"github.com/gorilla/websocket"
	"github.com/juju/errors"
)

type UserConnectionDao interface {
	Save(uint64, *websocket.Conn)
	FindByID(uint64) (*websocket.Conn, error)
	Delete(uint64)
}

type UserConnectionDaoImpl struct {
	users *hashmap.HashMap
}

var userConnectionDao UserConnectionDao

func init() {
	userConnectionDao = &UserConnectionDaoImpl{
		users: &hashmap.HashMap{},
	}
}

func NewUserConnectionDao() UserConnectionDao {
	return userConnectionDao
}

func (u *UserConnectionDaoImpl) Save(userID uint64, conn *websocket.Conn) {
	u.users.Set(userID, conn)
}

func (u *UserConnectionDaoImpl) FindByID(userID uint64) (*websocket.Conn, error) {
	conn, ok := u.users.Get(userID)
	if !ok {
		return nil, errors.NewNotFound(nil, fmt.Sprint(userID))
	}
	return conn.(*websocket.Conn), nil
}

func (u *UserConnectionDaoImpl) Delete(userID uint64) {
	u.users.Del(userID)
}
