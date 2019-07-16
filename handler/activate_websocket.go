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
		log.Printf("error: %+v", err)
		return
	}

	const writeWait = 10 * time.Second
	const pongWait = 10 * time.Second
	const pingPeriod = (pongWait * 9) / 10
	tick := time.NewTicker(pingPeriod)
	conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("close: %d(%d)", user.ID, code)
		tick.Stop()

		userRoom, err := h.userRoom.FindByUserID(user.ID)
		if err != nil && !errors.IsNotFound(err) {
			return errors.Trace(err)
		}
		if userRoom == nil {
			return nil
		}
		room, err := h.room.FindByID(userRoom.RoomID)
		if err != nil {
			return errors.Trace(err)
		}
		room.DelUser(user.ID)
		if len(room.ListUsers()) > 0 {
			room.Broadcast(&response.Stream{
				Type: response.StreamTypeLeave,
				User: user,
			})
		}
		return nil
	})

	go func() {
		defer conn.Close()

		for {
			select {
			case _, ok := <-tick.C:
				if !ok {
					return
				}
				if err := conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
					return
				}
				if err := conn.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
					return
				}
			}
		}
	}()

	go func() {
		defer conn.Close()

		if err := conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			return
		}

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				if err == websocket.ErrReadLimit {
					log.Printf("ping timeouted: %d", user.ID)
				} else {
					log.Printf("read err: %+v", err)
				}
				return
			}

			if string(msg) == "pong" {
				if err := conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
					return
				}
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
