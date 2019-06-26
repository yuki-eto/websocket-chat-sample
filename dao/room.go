package dao

import (
	"sync"
	"websocket-chat-sample/entity"

	"github.com/juju/errors"
)

type RoomDao interface {
	Create(*entity.Room) error
	FindByID(string) (*entity.Room, error)
}

type RoomDaoImpl struct {
	mtx   *sync.Mutex
	rooms map[string]*entity.Room
}

func NewRoomDao() RoomDao {
	return &RoomDaoImpl{
		mtx:   new(sync.Mutex),
		rooms: make(map[string]*entity.Room),
	}
}

func (r *RoomDaoImpl) Create(room *entity.Room) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if existRoom, _ := r.FindByID(room.ID); existRoom != nil {
		return errors.NewAlreadyExists(nil, room.ID)
	}

	r.rooms[room.ID] = room
	return nil
}

func (r *RoomDaoImpl) FindByID(id string) (*entity.Room, error) {
	return r.rooms[id], nil
}
