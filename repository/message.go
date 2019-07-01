package repository

import (
	"encoding/json"
	"websocket-chat-sample/dao"
	"websocket-chat-sample/entity"

	"github.com/juju/errors"
)

type MessageInstance struct {
	*entity.Message
}

func NewMessageInstance(message *entity.Message) *MessageInstance {
	return &MessageInstance{
		Message: message,
	}
}

type MessagesInstance struct {
	values []*MessageInstance
}

func NewMessagesInstance(entities []*entity.Message) *MessagesInstance {
	messages := new(MessagesInstance)
	for _, e := range entities {
		messages.Add(NewMessageInstance(e))
	}
	return messages
}
func (m *MessagesInstance) Clear() {
	m.values = []*MessageInstance{}
}
func (m *MessagesInstance) Add(msg *MessageInstance) {
	m.values = append(m.values, msg)
}
func (m *MessagesInstance) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(m.values)
	return b, errors.Trace(err)
}

type MessageRepository interface {
	Create(string) error
	Delete(string) error
	FindByRoomID(string) (*MessagesInstance, error)
	Push(*entity.Message) error
}

type MessageRepositoryImpl struct {
	messageDao dao.MessageDao
}

func NewMessageRepository() MessageRepository {
	return &MessageRepositoryImpl{
		messageDao: dao.NewMessageDao(),
	}
}

func (r *MessageRepositoryImpl) Create(id string) error {
	r.messageDao.Create(id)
	return nil
}
func (r *MessageRepositoryImpl) Delete(id string) error {
	return errors.Trace(r.messageDao.Delete(id))
}
func (r *MessageRepositoryImpl) FindByRoomID(id string) (*MessagesInstance, error) {
	msgs, err := r.messageDao.FindByRoomID(id)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return NewMessagesInstance(msgs), nil
}
func (r *MessageRepositoryImpl) Push(msg *entity.Message) error {
	return errors.Trace(r.messageDao.Push(msg))
}
