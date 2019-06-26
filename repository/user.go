package repository

import (
	"fmt"
	"websocket-chat-sample/dao"
	"websocket-chat-sample/entity"

	"github.com/juju/errors"
)

type UserInstance struct {
	*entity.User
}

func NewUserInstance(user *entity.User) *UserInstance {
	return &UserInstance{
		User: user,
	}
}

type UsersInstance struct {
	values map[uint64]*UserInstance
}

func NewUsersInstance() *UsersInstance {
	users := new(UsersInstance)
	users.Clear()
	return users
}
func (u *UsersInstance) Clear() {
	u.values = make(map[uint64]*UserInstance)
}
func (u *UsersInstance) Set(user *UserInstance) {
	u.values[user.ID] = user
}
func (u *UsersInstance) Delete(id uint64) {
	user := u.values[id]
	if user == nil {
		return
	}
	delete(u.values, id)
}
func (u *UsersInstance) FindByID(id uint64) *UserInstance {
	return u.values[id]
}
func (u *UsersInstance) Each(f func(instance *UserInstance)) {
	for _, value := range u.values {
		f(value)
	}
}
func (u *UsersInstance) EachWithError(f func(instance *UserInstance) error) error {
	for _, value := range u.values {
		if err := f(value); err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

type UserRepository interface {
	Create(string, string) (*UserInstance, error)
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

func (u *UserRepositoryImpl) Create(token, name string) (*UserInstance, error) {
	user := &entity.User{Token: token, Name: name}
	if err := u.userDao.Create(user); err != nil {
		return nil, errors.Trace(err)
	}
	return NewUserInstance(user), nil
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
