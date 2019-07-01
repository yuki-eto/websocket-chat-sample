package handler

import (
	"log"
	"net/http"
	"websocket-chat-sample/repository"
	"websocket-chat-sample/response"

	"github.com/gorilla/websocket"
)

type ActivateWebSocket struct {
	upgrader *websocket.Upgrader
	user     repository.UserRepository
}

func NewActivateWebsocketHandler(upgrader *websocket.Upgrader) http.Handler {
	return &ActivateWebSocket{
		upgrader: upgrader,
		user:     repository.NewUserRepository(),
	}
}

func (h *ActivateWebSocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := GetUser(w, r, h.user)
	if user == nil {
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		res := &Response{
			&response.ErrorResponse{Error: err},
		}
		res.InternalError(w)
		return
	}
	defer conn.Close()

	doneCh := make(chan bool, 1)
	conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("close: %d, %s", code, text)
		doneCh <- true
		return nil
	})

	if err := user.SetWSConnection(conn); err != nil {
		log.Printf("stream connection error: %+v", err)
		return
	}
	log.Printf("stream connected: %d", user.ID)

	<-doneCh
	log.Printf("stream disconnected: %d", user.ID)
}
