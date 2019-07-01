package dao

import (
	"sync"
	"websocket-chat-sample/entity"

	"github.com/juju/errors"
)

type RoomDao interface {
	Create(*entity.Room) error
	Update(*entity.Room) error
	Delete(string) error
	FindByID(string) (*entity.Room, error)
}

type RoomDaoImpl struct {
	rooms *sync.Map
}

var roomDao *RoomDaoImpl

func init() {
	roomDao = &RoomDaoImpl{
		rooms: new(sync.Map),
	}
}

func NewRoomDao() RoomDao {
	return roomDao
}

func (r *RoomDaoImpl) Create(room *entity.Room) error {
	if existRoom, _ := r.FindByID(room.ID); existRoom != nil {
		return errors.NewAlreadyExists(nil, room.ID)
	}

	r.rooms.Store(room.ID, room)
	return nil
}

func (r *RoomDaoImpl) Update(room *entity.Room) error {
	return nil
}

func (r *RoomDaoImpl) Delete(roomID string) error {
	r.rooms.Delete(roomID)
	return nil
}

func (r *RoomDaoImpl) FindByID(id string) (*entity.Room, error) {
	room, ok := r.rooms.Load(id)
	if !ok {
		return nil, nil
	}
	return room.(*entity.Room), nil
}
