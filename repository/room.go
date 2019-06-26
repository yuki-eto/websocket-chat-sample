package repository

import (
	"github.com/juju/errors"
	"websocket-chat-sample/dao"
	"websocket-chat-sample/entity"
)

type RoomInstance struct {
	*entity.Room
}
func NewRoomInstance(room *entity.Room) *RoomInstance {
	return &RoomInstance{
		Room: room,
	}
}

type RoomRepository interface {
	Create(string, string) (*RoomInstance, error)
	FindByID(string) (*RoomInstance, error)
}

type RoomRepositoryImpl struct {
	roomDao dao.RoomDao
}

func NewRoomRepository() RoomRepository {
	return &RoomRepositoryImpl{
		roomDao: dao.NewRoomDao(),
	}
}

func (r *RoomRepositoryImpl) Create(id, name string) (*RoomInstance, error) {
	room := &entity.Room{ID: id, Name: name}
	if err := r.roomDao.Create(room); err != nil {
		return nil, errors.Trace(err)
	}
	return NewRoomInstance(room), nil
}

func (r *RoomRepositoryImpl) FindByID(id string) (*RoomInstance, error) {
	room, err := r.roomDao.FindByID(id)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if room == nil {
		return nil, errors.NewNotFound(nil, id)
	}
	return NewRoomInstance(room), nil
}
