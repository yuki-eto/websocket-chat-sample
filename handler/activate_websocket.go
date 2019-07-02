package handler

import (
	"log"
	"net/http"
	"time"
	"websocket-chat-sample/repository"
	"websocket-chat-sample/response"

	"github.com/gorilla/websocket"
	"github.com/juju/errors"
)

type ActivateWebSocket struct {
	upgrader *websocket.Upgrader
	user     repository.UserRepository
	room     repository.RoomRepository
	userRoom repository.UserRoomRepository
}

func NewActivateWebsocketHandler(upgrader *websocket.Upgrader) http.Handler {
	return &ActivateWebSocket{
		upgrader: upgrader,
		user:     repository.NewUserRepository(),
		room:     repository.NewRoomRepository(),
		userRoom: repository.NewUserRoomRepository(),
	}
}

func (h *ActivateWebSocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := GetUser(w, r, h.user)
	if user == nil {
		return
	}

	log.Printf("upgrade to websocket: %s", r.RemoteAddr)

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		res := &Response{
			&response.ErrorResponse{Error: err},
		}
		res.InternalError(w)
		return
	}

	tick := time.NewTicker(time.Second * 10)
	conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("close: %d, %s", code, text)
		tick.Stop()
		return nil
	})

	go func() {
		for {
			select {
			case _, ok := <-tick.C:
				if !ok {
					return
				}
				if err := conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
					log.Printf("ping error: %+v", err)
					conn.Close()
					return
				}
			default:
				return
			}
		}
	}()

	user.SetWSConnection(conn)
	if err := h.user.SaveConnection(user); err != nil {
		log.Printf("save connection error: %+v", err)
		return
	}

	userRoom, err := h.userRoom.FindByUserID(user.ID)
	if err != nil && !errors.IsNotFound(err) {
		res := &Response{
			&response.ErrorResponse{Error: err},
		}
		res.InternalError(w)
		conn.Close()
		return
	}

	if userRoom != nil {
		room, err := h.room.FindByID(userRoom.RoomID)
		if err != nil {
			res := &Response{
				&response.ErrorResponse{Error: err},
			}
			res.InternalError(w)
			conn.Close()
			return
		}
		room.SetUser(user)
	}

	log.Printf("stream connected: %d", user.ID)
}
