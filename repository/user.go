package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"websocket-chat-sample/dao"
	"websocket-chat-sample/entity"

	"github.com/cornelk/hashmap"
	"github.com/gorilla/websocket"
	"github.com/juju/errors"
)

type UserInstance struct {
	*entity.User

	conn *websocket.Conn
}

func NewUserInstance(user *entity.User, conn *websocket.Conn) *UserInstance {
	return &UserInstance{
		User: user,
		conn: conn,
	}
}

func (u *UserInstance) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(u.User)
	return b, errors.Trace(err)
}

func (u *UserInstance) SetWSConnection(conn *websocket.Conn) {
	if u.conn != nil {
		u.conn.Close()
	}
	u.conn = conn
}

type UsersInstance struct {
	values *hashmap.HashMap
}

func NewUsersInstance() *UsersInstance {
	users := new(UsersInstance)
	users.Clear()
	return users
}
func (u *UsersInstance) Clear() {
	u.values = &hashmap.HashMap{}
}
func (u *UsersInstance) Set(user *UserInstance) {
	u.values.Set(user.ID, user)
}
func (u *UsersInstance) Delete(id uint64) {
	u.values.Del(id)
}
func (u *UsersInstance) FindByID(id uint64) *UserInstance {
	user, ok := u.values.Get(id)
	if !ok {
		return nil
	}
	return user.(*UserInstance)
}
func (u *UsersInstance) Each(f func(instance *UserInstance)) {
	for item := range u.values.Iter() {
		f(item.Value.(*UserInstance))
	}
}
func (u *UsersInstance) EachWithError(f func(instance *UserInstance) error) error {
	for item := range u.values.Iter() {
		if err := f(item.Value.(*UserInstance)); err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}
func (u *UsersInstance) Broadcast(msg interface{}) {
	u.Each(func(instance *UserInstance) {
		if instance.conn == nil {
			log.Printf("connection is null: %d", instance.ID)
			return
		}
		if err := instance.conn.WriteJSON(msg); err != nil {
			log.Printf("broadcast err: %+v", err)
		}
	})
}
func (u *UsersInstance) List() (list []*UserInstance) {
	u.Each(func(instance *UserInstance) {
		list = append(list, instance)
	})
	return list
}

type UserRepository interface {
	Create(*UserInstance) error
	Save(*UserInstance) error
	SaveConnection(*UserInstance) error
	FindByID(uint64) (*UserInstance, error)
	FindByToken(string) (*UserInstance, error)
}

type UserRepositoryImpl struct {
	userDao           dao.UserDao
	userConnectionDao dao.UserConnectionDao
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{
		userDao:           dao.NewUserDao(),
		userConnectionDao: dao.NewUserConnectionDao(),
	}
}

func (u *UserRepositoryImpl) Create(user *UserInstance) error {
	return errors.Trace(u.userDao.Create(user.User))
}

func (u *UserRepositoryImpl) Save(user *UserInstance) error {
	if user.ID == 0 {
		return errors.Trace(u.Create(user))
	}
	return errors.Trace(u.userDao.Update(user.User))
}

func (u *UserRepositoryImpl) SaveConnection(user *UserInstance) error {
	if user.conn == nil {
		return nil
	}
	u.userConnectionDao.Save(user.ID, user.conn)
	return nil
}

func (u *UserRepositoryImpl) FindByID(id uint64) (*UserInstance, error) {
	user, err := u.userDao.FindByID(id)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if user == nil {
		return nil, errors.NewNotFound(nil, fmt.Sprint(id))
	}
	conn, err := u.userConnectionDao.FindByID(id)
	if err != nil && !errors.IsNotFound(err) {
		return nil, errors.Trace(err)
	}
	return NewUserInstance(user, conn), nil
}

func (u *UserRepositoryImpl) FindByToken(token string) (*UserInstance, error) {
	user, err := u.userDao.FindByToken(token)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if user == nil {
		return nil, errors.NewNotFound(nil, token)
	}
	conn, err := u.userConnectionDao.FindByID(user.ID)
	if err != nil && !errors.IsNotFound(err) {
		return nil, errors.Trace(err)
	}
	return NewUserInstance(user, conn), nil
}
