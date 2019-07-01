package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"websocket-chat-sample/repository"
	"websocket-chat-sample/request"
	"websocket-chat-sample/response"

	"github.com/juju/errors"
)

type JoinRoom struct {
	user     repository.UserRepository
	room     repository.RoomRepository
	userRoom repository.UserRoomRepository
	message  repository.MessageRepository
}

func NewJoinRoomHandler() http.Handler {
	return &JoinRoom{
		user:     repository.NewUserRepository(),
		room:     repository.NewRoomRepository(),
		userRoom: repository.NewUserRoomRepository(),
		message:  repository.NewMessageRepository(),
	}
}

func (h *JoinRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := GetUser(w, r, h.user)
	if user == nil {
		return
	}

	userRoom, err := h.userRoom.FindByUserID(user.ID)
	if err != nil && !errors.IsNotFound(err) {
		res := &Response{
			&response.ErrorResponse{Error: err},
		}
		res.InternalError(w)
		return
	}
	if userRoom != nil {
		res := &Response{
			&response.ErrorResponse{
				Error: errors.NewBadRequest(nil, fmt.Sprintf("already in room : %s", userRoom.RoomID)),
			},
		}
		res.BadRequest(w)
		return
	}

	var req request.JoinRoom
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		res := &Response{
			&response.ErrorResponse{Error: err},
		}
		res.BadRequest(w)
		return
	}
	roomID := req.RoomID

	room, err := h.room.FindByID(roomID)
	if err != nil {
		res := &Response{
			&response.ErrorResponse{Error: err},
		}
		res.InternalError(w)
		return
	}
	messages, err := h.message.FindByRoomID(roomID)
	if err != nil {
		res := &Response{
			&response.ErrorResponse{Error: err},
		}
		res.InternalError(w)
		return
	}

	room.AddUser(user)

	msg := &response.Stream{
		Type: response.StreamTypeJoin,
		User: user,
	}
	room.Broadcast(msg)

	res := &Response{
		&response.JoinRoom{
			Room:     room,
			Users:    room.ListUsers(),
			Messages: messages,
		},
	}
	res.Ok(w)
}
