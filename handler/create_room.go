package handler

import (
	"encoding/json"
	"net/http"
	"websocket-chat-sample/entity"
	"websocket-chat-sample/repository"
	"websocket-chat-sample/request"
	"websocket-chat-sample/response"

	"github.com/juju/errors"
)

type CreateRoom struct {
	user    repository.UserRepository
	room    repository.RoomRepository
	message repository.MessageRepository
}

func NewCreateRoomHandler() http.Handler {
	return &CreateRoom{
		user:    repository.NewUserRepository(),
		room:    repository.NewRoomRepository(),
		message: repository.NewMessageRepository(),
	}
}

func (h *CreateRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := GetUser(w, r, h.user)
	if user == nil {
		return
	}

	var req request.CreateRoom
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		res := &Response{
			&response.ErrorResponse{Error: err},
		}
		res.BadRequest(w)
		return
	}

	if err := h.message.Create(req.RoomID); err != nil {
		res := &Response{
			&response.ErrorResponse{Error: err},
		}
		res.InternalError(w)
		return
	}

	room := repository.NewRoomInstance(&entity.Room{
		ID:   req.RoomID,
		Name: "",
	})
	if err := h.room.Create(room); err != nil {
		errBody := &response.ErrorResponse{Error: err}
		res := &Response{errBody}
		if errors.IsAlreadyExists(err) {
			errBody.Code = response.ErrorCodeAlreadyExistRoom
			res.BadRequest(w)
		} else {
			res.InternalError(w)
		}
		return
	}

	res := &Response{}
	res.NoContent(w)
}
