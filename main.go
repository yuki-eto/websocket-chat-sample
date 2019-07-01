package main

import (
	"fmt"
	"log"
	"net/http"
	"websocket-chat-sample/handler"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func main() {
	router := mux.NewRouter()
	upgrader := &websocket.Upgrader{}

	router.Handle("/create_user", handler.NewCreateUserHandler()).Methods(http.MethodPost)
	router.Handle("/login", handler.NewLoginHandler()).Methods(http.MethodPost)
	router.Handle("/create_room", handler.NewCreateRoomHandler()).Methods(http.MethodPost)
	router.Handle("/join_room", handler.NewJoinRoomHandler()).Methods(http.MethodPost)
	router.Handle("/leave_room", handler.NewLeaveRoomHandler()).Methods(http.MethodPost)
	router.Handle("/message", handler.NewMessageHandler()).Methods(http.MethodPost)

	router.Handle("/websocket", handler.NewActivateWebsocketHandler(upgrader))

	addr := fmt.Sprintf("%s:%d", "0.0.0.0", 19999)
	log.Printf("start server on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		panic(err)
	}
}
