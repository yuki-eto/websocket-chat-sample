package repository

import (
	"encoding/json"
	"websocket-chat-sample/dao"
	"websocket-chat-sample/entity"

	"github.com/cornelk/hashmap"
	"github.com/juju/errors"
)

type RoomInstance struct {
	*entity.Room
	users     *UsersInstance
	isCreated bool
}

func NewRoomInstance(room *entity.Room) *RoomInstance {
	return &RoomInstance{
		Room:      room,
		users:     NewUsersInstance(),
		isCreated: false,
	}
}

func (r *RoomInstance) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(r.Room)
	return b, errors.Trace(err)
}

func (r *RoomInstance) SetUser(user *UserInstance) {
	r.users.Set(user)
}
func (r *RoomInstance) DelUser(userID uint64) {
	r.users.Delete(userID)
}
func (r *RoomInstance) ListUsers() []*UserInstance {
	return r.users.List()
}
func (r *RoomInstance) Broadcast(msg interface{}) {
	r.users.Broadcast(msg)
}
func (r *RoomInstance) Close() {
	// TODO: ルーム閉じる時に必要な処理を入れる
}

type RoomRepository interface {
	Create(*RoomInstance) error
	Save(*RoomInstance) error
	Delete(*RoomInstance) error
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

func (r *RoomRepositoryImpl) Create(room *RoomInstance) error {
	if err := r.roomDao.Create(room.Room); err != nil {
		return errors.Trace(err)
	}
	room.isCreated = true
	roomCache.Set(room)
	return nil
}

func (r *RoomRepositoryImpl) Save(room *RoomInstance) error {
	if !room.isCreated {
		return errors.Trace(r.Create(room))
	}
	return errors.Trace(r.roomDao.Update(room.Room))
}

func (r *RoomRepositoryImpl) Delete(room *RoomInstance) error {
	room.Close()
	return errors.Trace(r.roomDao.Delete(room.ID))
}

func (r *RoomRepositoryImpl) FindByID(id string) (*RoomInstance, error) {
	roomInstance := roomCache.Get(id)
	if roomInstance != nil {
		return roomInstance, nil
	}

	room, err := r.roomDao.FindByID(id)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if room == nil {
		return nil, errors.NewNotFound(nil, id)
	}
	instance := NewRoomInstance(room)
	instance.isCreated = true
	roomCache.Set(instance)
	return instance, nil
}

type roomCaches struct {
	list *hashmap.HashMap
}

func (r *roomCaches) Set(room *RoomInstance) {
	r.list.Set(room.ID, room)
}
func (r *roomCaches) Get(id string) *RoomInstance {
	room, ok := r.list.Get(id)
	if !ok {
		return nil
	}
	return room.(*RoomInstance)
}

var roomCache *roomCaches

func init() {
	roomCache = &roomCaches{
		list: &hashmap.HashMap{},
	}
}
