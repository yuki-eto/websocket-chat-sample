package repository

import (
	"fmt"
	"websocket-chat-sample/dao"
	"websocket-chat-sample/entity"

	"github.com/juju/errors"
)

type UserRoomInstance struct {
	*entity.UserRoom

	isCreated bool
}

func NewUserRoomInstance(uRoom *entity.UserRoom) *UserRoomInstance {
	return &UserRoomInstance{
		UserRoom:  uRoom,
		isCreated: false,
	}
}

type UserRoomRepository interface {
	Create(*UserRoomInstance) error
	Save(*UserRoomInstance) error
	Delete(*UserRoomInstance) error
	FindByUserID(uint64) (*UserRoomInstance, error)
}

type UserRoomRepositoryImpl struct {
	userRoomDao dao.UserRoomDao
}

func NewUserRoomRepository() UserRoomRepository {
	return &UserRoomRepositoryImpl{
		userRoomDao: dao.NewUserRoomDao(),
	}
}

func (u *UserRoomRepositoryImpl) Create(uRoom *UserRoomInstance) error {
	if err := u.userRoomDao.Create(uRoom.UserRoom); err != nil {
		return errors.Trace(err)
	}
	uRoom.isCreated = true
	return nil
}

func (u *UserRoomRepositoryImpl) Save(uRoom *UserRoomInstance) error {
	if !uRoom.isCreated {
		return errors.Trace(u.Create(uRoom))
	}
	return errors.Trace(u.userRoomDao.Update(uRoom.UserRoom))
}

func (u *UserRoomRepositoryImpl) Delete(uRoom *UserRoomInstance) error {
	return errors.Trace(u.userRoomDao.Delete(uRoom.UserID))
}

func (u *UserRoomRepositoryImpl) FindByUserID(userID uint64) (*UserRoomInstance, error) {
	uRoom, err := u.userRoomDao.FindByUserID(userID)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if uRoom == nil {
		return nil, errors.NewNotFound(nil, fmt.Sprint(userID))
	}
	instance := NewUserRoomInstance(uRoom)
	instance.isCreated = true
	return instance, nil
}
