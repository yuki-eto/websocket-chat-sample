package dao

import (
	"sync"
	"websocket-chat-sample/entity"

	"github.com/juju/errors"
)

type UserRoomDao interface {
	Create(*entity.UserRoom) error
	Update(*entity.UserRoom) error
	Delete(uint64) error
	FindByUserID(uint64) (*entity.UserRoom, error)
}

type UserRoomDaoImpl struct {
	userRooms *sync.Map
}

var userRoomDao *UserRoomDaoImpl

func init() {
	userRoomDao = &UserRoomDaoImpl{
		userRooms: new(sync.Map),
	}
}

func NewUserRoomDao() UserRoomDao {
	return userRoomDao
}

func (u *UserRoomDaoImpl) Create(userRoom *entity.UserRoom) error {
	if user, err := u.FindByUserID(userRoom.UserID); err != nil {
		return errors.Trace(err)
	} else if user != nil {
		return errors.NewAlreadyExists(nil, "already created user_room")
	}

	u.userRooms.Store(userRoom.UserID, userRoom)
	return nil
}

func (u *UserRoomDaoImpl) Update(userRoom *entity.UserRoom) error {
	return nil
}

func (u *UserRoomDaoImpl) Delete(userID uint64) error {
	u.userRooms.Delete(userID)
	return nil
}

func (u *UserRoomDaoImpl) FindByUserID(userID uint64) (*entity.UserRoom, error) {
	uRoom, ok := u.userRooms.Load(userID)
	if !ok {
		return nil, nil
	}
	return uRoom.(*entity.UserRoom), nil
}
