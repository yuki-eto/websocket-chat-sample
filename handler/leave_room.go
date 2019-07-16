package handler

import (
	"net/http"
	"websocket-chat-sample/repository"
	"websocket-chat-sample/response"

	"github.com/juju/errors"
)

type LeaveRoom struct {
	user     repository.UserRepository
	room     repository.RoomRepository
	userRoom repository.UserRoomRepository
	message  repository.MessageRepository
}

func NewLeaveRoomHandler() http.Handler {
	return &LeaveRoom{
		user:     repository.NewUserRepository(),
		room:     repository.NewRoomRepository(),
		userRoom: repository.NewUserRoomRepository(),
		message:  repository.NewMessageRepository(),
	}
}

func (h *LeaveRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := GetUser(w, r, h.user)
	if user == nil {
		return
	}

	userRoom, err := h.userRoom.FindByUserID(user.ID)
	if err != nil {
		res := &Response{
			&response.ErrorResponse{Error: err},
		}
		if errors.IsNotFound(err) {
			res.NotFound(w)
		} else {
			res.InternalError(w)
		}
		return
	}

	roomID := userRoom.RoomID

	room, err := h.room.FindByID(roomID)
	if err != nil && !errors.IsNotFound(err) {
		res := &Response{
			&response.ErrorResponse{Error: err},
		}
		res.InternalError(w)
		return
	}

	if err := h.userRoom.Delete(userRoom); err != nil {
		res := &Response{
			&response.ErrorResponse{Error: err},
		}
		res.InternalError(w)
		return
	}

	if room != nil {
		if err := h.LeaveRoom(room, user); err != nil {
			res := &Response{
				&response.ErrorResponse{Error: err},
			}
			res.InternalError(w)
			return
		}
	}

	res := &Response{}
	res.NoContent(w)
}

func (h *LeaveRoom) LeaveRoom(room *repository.RoomInstance, user *repository.UserInstance) error {
	room.DelUser(user.ID)
	list := room.ListUsers()

	if len(list) > 0 {
		room.Broadcast(&response.Stream{
			Type: response.StreamTypeLeave,
			User: user,
		})
	}

	return nil
}
