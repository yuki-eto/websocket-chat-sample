package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"websocket-chat-sample/entity"
	"websocket-chat-sample/repository"
	"websocket-chat-sample/request"
	"websocket-chat-sample/response"

	"github.com/juju/errors"
)

type Message struct {
	user     repository.UserRepository
	room     repository.RoomRepository
	userRoom repository.UserRoomRepository
	message  repository.MessageRepository
}

func NewMessageHandler() http.Handler {
	return &Message{
		user:     repository.NewUserRepository(),
		room:     repository.NewRoomRepository(),
		userRoom: repository.NewUserRoomRepository(),
		message:  repository.NewMessageRepository(),
	}
}

func (h *Message) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	var req request.Message
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		res := &Response{
			&response.ErrorResponse{Error: err},
		}
		res.BadRequest(w)
		return
	}

	room, err := h.room.FindByID(roomID)
	if err != nil {
		res := &Response{
			&response.ErrorResponse{Error: err},
		}
		res.InternalError(w)
		return
	}

	now := time.Now()
	msg := &entity.Message{
		RoomID: roomID,
		UserID: user.ID,
		Text:   req.Text,
		Time:   &now,
	}
	if err := h.message.Push(msg); err != nil {
		res := &Response{
			&response.ErrorResponse{Error: err},
		}
		res.InternalError(w)
		return
	}

	stream := &response.Stream{
		Type:    response.StreamTypeChat,
		Message: repository.NewMessageInstance(msg),
	}
	room.Broadcast(stream)

	res := &Response{}
	res.NoContent(w)
}
