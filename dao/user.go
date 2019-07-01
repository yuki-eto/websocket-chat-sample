package dao

import (
	"sync"
	"websocket-chat-sample/entity"

	"github.com/juju/errors"
)

type UserDao interface {
	Create(*entity.User) error
	Update(*entity.User) error
	FindByID(uint64) (*entity.User, error)
	FindByToken(string) (*entity.User, error)
}

type UserDaoImpl struct {
	users       *sync.Map
	userByToken *sync.Map

	sequence    uint64
	sequenceMtx *sync.Mutex
}

var userDao *UserDaoImpl

func init() {
	userDao = &UserDaoImpl{
		users:       new(sync.Map),
		userByToken: new(sync.Map),

		sequence:    0,
		sequenceMtx: new(sync.Mutex),
	}
}

func NewUserDao() UserDao {
	return userDao
}

func (u *UserDaoImpl) sequenceID() uint64 {
	u.sequenceMtx.Lock()
	u.sequence++
	u.sequenceMtx.Unlock()
	return u.sequence
}

func (u *UserDaoImpl) Create(user *entity.User) error {
	existUser, _ := u.FindByToken(user.LoginToken)
	if existUser != nil {
		return errors.NewAlreadyExists(nil, user.LoginToken)
	}

	user.ID = u.sequenceID()
	u.users.Store(user.ID, user)
	u.userByToken.Store(user.LoginToken, user)
	return nil
}

func (u *UserDaoImpl) Update(user *entity.User) error {
	return nil
}

func (u *UserDaoImpl) FindByID(id uint64) (*entity.User, error) {
	user, ok := u.users.Load(id)
	if !ok {
		return nil, nil
	}
	return user.(*entity.User), nil
}

func (u *UserDaoImpl) FindByToken(token string) (*entity.User, error) {
	user, ok := u.userByToken.Load(token)
	if !ok {
		return nil, nil
	}
	return user.(*entity.User), nil
}
