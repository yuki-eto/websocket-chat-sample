package repository

import (
	"testing"
	"websocket-chat-sample/entity"

	"github.com/juju/errors"
)

func TestUsersInstance_Each(t *testing.T) {
	users := NewUsersInstance()
	users.Set(&UserInstance{User: &entity.User{ID: 1}})
	users.Set(&UserInstance{User: &entity.User{ID: 2}})
	users.Set(&UserInstance{User: &entity.User{ID: 3}})
	users.Set(&UserInstance{User: &entity.User{ID: 4}})
	users.Set(&UserInstance{User: &entity.User{ID: 5}})

	i := uint64(0)
	users.Each(func(instance *UserInstance) {
		i += instance.ID
	})
	if i != 15 {
		t.Fatalf("total ids not match: %+v <> %+v", 15, i)
	}
}

func TestUsersInstance_EachWithError(t *testing.T) {
	users := NewUsersInstance()
	users.Set(&UserInstance{User: &entity.User{ID: 1}})
	users.Set(&UserInstance{User: &entity.User{ID: 2}})
	users.Set(&UserInstance{User: &entity.User{ID: 3}})
	users.Set(&UserInstance{User: &entity.User{ID: 4}})
	users.Set(&UserInstance{User: &entity.User{ID: 5}})

	i := uint64(0)
	if err := users.EachWithError(func(instance *UserInstance) error {
		if instance.ID == 3 {
			return errors.New("error")
		}
		return nil
	}); err != nil {
		if err.Error() != "error" {
			t.Fatalf("does not match error: %s", err.Error())
		}
	}

	if i == 15 {
		t.Fatalf("total ids not match: %+v <> %+v", 15, i)
	}
}
