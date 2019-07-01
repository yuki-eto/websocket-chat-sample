package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"websocket-chat-sample/dao"
	"websocket-chat-sample/entity"

	"github.com/gorilla/websocket"
	"github.com/juju/errors"
)

type UserInstance struct {
	*entity.User

	mtx  *sync.Mutex
	conn *websocket.Conn
}

func NewUserInstance(user *entity.User) *UserInstance {
	return &UserInstance{
		User: user,
		mtx:  new(sync.Mutex),
		conn: nil,
	}
}

func (u *UserInstance) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(u.User)
	return b, errors.Trace(err)
}

func (u *UserInstance) SetWSConnection(conn *websocket.Conn) error {
	u.mtx.Lock()
	defer u.mtx.Unlock()

	if u.conn != nil {
		if err := u.conn.Close(); err != nil {
			return errors.Trace(err)
		}
	}

	u.conn = conn
	return nil
}

type UsersInstance struct {
	values *sync.Map
}

func NewUsersInstance() *UsersInstance {
	users := new(UsersInstance)
	users.Clear()
	return users
}
func (u *UsersInstance) Clear() {
	u.values = new(sync.Map)
}
func (u *UsersInstance) Set(user *UserInstance) {
	u.values.Store(user.ID, user)
}
func (u *UsersInstance) Delete(id uint64) {
	u.values.Delete(id)
}
func (u *UsersInstance) FindByID(id uint64) *UserInstance {
	user, ok := u.values.Load(id)
	if !ok {
		return nil
	}
	return user.(*UserInstance)
}
func (u *UsersInstance) Each(f func(instance *UserInstance)) {
	u.values.Range(func(key, value interface{}) bool {
		user := value.(*UserInstance)
		f(user)
		return true
	})
}
func (u *UsersInstance) EachWithError(f func(instance *UserInstance) error) error {
	var err error
	u.values.Range(func(key, value interface{}) bool {
		user := value.(*UserInstance)
		if err = f(user); err != nil {
			return false
		}
		return true
	})
	return errors.Trace(err)
}
func (u *UsersInstance) Broadcast(msg interface{}) {
	u.Each(func(instance *UserInstance) {
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
	FindByID(uint64) (*UserInstance, error)
	FindByToken(string) (*UserInstance, error)
}

type UserRepositoryImpl struct {
	userDao dao.UserDao
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{
		userDao: dao.NewUserDao(),
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

func (u *UserRepositoryImpl) FindByID(id uint64) (*UserInstance, error) {
	user, err := u.userDao.FindByID(id)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if user == nil {
		return nil, errors.NewNotFound(nil, fmt.Sprint(id))
	}
	return NewUserInstance(user), nil
}

func (u *UserRepositoryImpl) FindByToken(token string) (*UserInstance, error) {
	user, err := u.userDao.FindByToken(token)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if user == nil {
		return nil, errors.NewNotFound(nil, token)
	}
	return NewUserInstance(user), nil
}
