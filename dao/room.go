package dao

import (
	"websocket-chat-sample/entity"

	"github.com/cornelk/hashmap"
	"github.com/juju/errors"
)

type RoomDao interface {
	Create(*entity.Room) error
	Update(*entity.Room) error
	Delete(string) error
	FindByID(string) (*entity.Room, error)
}

type RoomDaoImpl struct {
	rooms *hashmap.HashMap
}

var roomDao *RoomDaoImpl

func init() {
	roomDao = &RoomDaoImpl{
		rooms: &hashmap.HashMap{},
	}
}

func NewRoomDao() RoomDao {
	return roomDao
}

func (r *RoomDaoImpl) Create(room *entity.Room) error {
	if existRoom, _ := r.FindByID(room.ID); existRoom != nil {
		return errors.NewAlreadyExists(nil, room.ID)
	}

	r.rooms.Set(room.ID, room)
	return nil
}

func (r *RoomDaoImpl) Update(room *entity.Room) error {
	return nil
}

func (r *RoomDaoImpl) Delete(roomID string) error {
	r.rooms.Del(roomID)
	return nil
}

func (r *RoomDaoImpl) FindByID(id string) (*entity.Room, error) {
	room, ok := r.rooms.Get(id)
	if !ok {
		return nil, nil
	}
	return room.(*entity.Room), nil
}
