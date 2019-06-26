package dao

import (
	"container/list"
	"sync"
	"websocket-chat-sample/entity"

	"github.com/juju/errors"
)

type MessagesStore struct {
	mtx         *sync.Mutex
	storeByRoom map[string]*list.List
	mtxByRoom   map[string]*sync.Mutex
}

func NewMessageStore() *MessagesStore {
	return &MessagesStore{
		mtx:         new(sync.Mutex),
		storeByRoom: make(map[string]*list.List),
		mtxByRoom:   make(map[string]*sync.Mutex),
	}
}
func (m *MessagesStore) New(roomID string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if _, exists := m.storeByRoom[roomID]; !exists {
		m.storeByRoom[roomID] = list.New()
	}
	if _, exists := m.mtxByRoom[roomID]; !exists {
		m.mtxByRoom[roomID] = new(sync.Mutex)
	}
}
func (m *MessagesStore) Push(msg *entity.Message) error {
	roomID := msg.RoomID
	store, exists := m.storeByRoom[roomID]
	if !exists {
		return errors.New("cannot find list")
	}
	mtx, exists := m.mtxByRoom[roomID]
	if !exists {
		return errors.New("cannot find mutex")
	}

	mtx.Lock()
	store.PushBack(msg)
	mtx.Unlock()

	return nil
}
func (m *MessagesStore) List(roomID string) (list []*entity.Message, err error) {
	store, exists := m.storeByRoom[roomID]
	if !exists {
		return list, errors.NewNotFound(nil, roomID)
	}

	for elm := store.Front(); elm != nil; elm = elm.Next() {
		msg := elm.Value.(*entity.Message)
		list = append(list, msg)
	}
	return list, nil
}

type MessageDao interface {
	Create(string)
	Push(*entity.Message) error
	FindByRoomID(string) ([]*entity.Message, error)
}

type MessageDaoImpl struct {
	store *MessagesStore
}

func NewMessageDao() MessageDao {
	return &MessageDaoImpl{
		store: NewMessageStore(),
	}
}

func (m *MessageDaoImpl) Create(roomID string) {
	m.store.New(roomID)
}
func (m *MessageDaoImpl) Push(msg *entity.Message) error {
	return errors.Trace(m.store.Push(msg))
}
func (m *MessageDaoImpl) FindByRoomID(roomID string) ([]*entity.Message, error) {
	msgs, err := m.store.List(roomID)
	return msgs, errors.Trace(err)
}
