package dao

import (
	"fmt"
	"sync"
	"websocket-chat-sample/entity"

	"github.com/juju/errors"
)

type UserDao interface {
	Create(*entity.User) error
	FindByID(uint64) (*entity.User, error)
	FindByToken(string) (*entity.User, error)
}

type UserDaoImpl struct {
	mtx         *sync.Mutex
	users       map[uint64]*entity.User
	userByToken map[string]*entity.User

	sequence    uint64
	sequenceMtx *sync.Mutex
}

func NewUserDao() UserDao {
	return &UserDaoImpl{
		mtx:         new(sync.Mutex),
		users:       make(map[uint64]*entity.User),
		userByToken: make(map[string]*entity.User),

		sequence:    0,
		sequenceMtx: new(sync.Mutex),
	}
}

func (u *UserDaoImpl) sequenceID() uint64 {
	u.sequenceMtx.Lock()
	u.sequence++
	u.sequenceMtx.Unlock()
	return u.sequence
}

func (u *UserDaoImpl) Create(user *entity.User) error {
	u.mtx.Lock()
	defer u.mtx.Unlock()

	existUser, _ := u.FindByToken(user.Token)
	if existUser != nil {
		return errors.NewAlreadyExists(nil, fmt.Sprint(user.Token))
	}

	user.ID = u.sequenceID()
	u.users[user.ID] = user
	u.userByToken[user.Token] = user
	return nil
}

func (u *UserDaoImpl) FindByID(id uint64) (*entity.User, error) {
	user := u.users[id]
	return user, nil
}

func (u *UserDaoImpl) FindByToken(token string) (*entity.User, error) {
	user := u.userByToken[token]
	return user, nil
}
