package dao

import (
	"container/list"
	"sync"
	"websocket-chat-sample/entity"

	"github.com/cornelk/hashmap"
	"github.com/juju/errors"
)

type MessagesStore struct {
	storeByRoom *hashmap.HashMap
	mtxByRoom   *hashmap.HashMap
}

func NewMessageStore() *MessagesStore {
	return &MessagesStore{
		storeByRoom: &hashmap.HashMap{},
		mtxByRoom:   &hashmap.HashMap{},
	}
}
func (m *MessagesStore) New(roomID string) {
	if _, exists := m.storeByRoom.Get(roomID); !exists {
		m.storeByRoom.Set(roomID, list.New())
	}
	if _, exists := m.mtxByRoom.Get(roomID); !exists {
		m.mtxByRoom.Set(roomID, new(sync.RWMutex))
	}
}
func (m *MessagesStore) getList(roomID string) *list.List {
	l, ok := m.storeByRoom.Get(roomID)
	if !ok {
		return nil
	}
	return l.(*list.List)
}
func (m *MessagesStore) getMutex(roomID string) *sync.RWMutex {
	mtx, ok := m.mtxByRoom.Get(roomID)
	if !ok {
		return nil
	}
	return mtx.(*sync.RWMutex)
}
func (m *MessagesStore) Push(msg *entity.Message) error {
	roomID := msg.RoomID

	store := m.getList(roomID)
	if store == nil {
		return errors.NewNotFound(nil, "cannot find store")
	}
	mtx := m.getMutex(roomID)
	if mtx == nil {
		return errors.NewNotFound(nil, "cannot find mutex")
	}

	mtx.Lock()
	store.PushBack(msg)
	mtx.Unlock()

	return nil
}
func (m *MessagesStore) List(roomID string) (list []*entity.Message, err error) {
	store := m.getList(roomID)
	if store == nil {
		return list, errors.NewNotFound(nil, "cannot find store")
	}
	mtx := m.getMutex(roomID)
	if mtx == nil {
		return nil, errors.NewNotFound(nil, "cannot find mutex")
	}

	mtx.RLock()
	for elm := store.Back(); elm != nil; elm = elm.Prev() {
		msg := elm.Value.(*entity.Message)
		list = append(list, msg)
	}
	mtx.RUnlock()

	return list, nil
}
func (m *MessagesStore) Delete(roomID string) error {
	m.storeByRoom.Del(roomID)
	m.mtxByRoom.Del(roomID)
	return nil
}

type MessageDao interface {
	Create(string)
	Push(*entity.Message) error
	Delete(string) error
	FindByRoomID(string) ([]*entity.Message, error)
}

type MessageDaoImpl struct {
	store *MessagesStore
}

var messageDao *MessageDaoImpl

func init() {
	messageDao = &MessageDaoImpl{
		store: NewMessageStore(),
	}
}

func NewMessageDao() MessageDao {
	return messageDao
}

func (m *MessageDaoImpl) Create(roomID string) {
	m.store.New(roomID)
}
func (m *MessageDaoImpl) Delete(roomID string) error {
	return errors.Trace(m.store.Delete(roomID))
}
func (m *MessageDaoImpl) Push(msg *entity.Message) error {
	return errors.Trace(m.store.Push(msg))
}
func (m *MessageDaoImpl) FindByRoomID(roomID string) ([]*entity.Message, error) {
	msgs, err := m.store.List(roomID)
	return msgs, errors.Trace(err)
}
