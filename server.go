package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"websocket-chat-sample/handler"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	router   = mux.NewRouter()
	upgrader = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main() {
	handleFunc("/create_user", handler.NewCreateUserHandler()).Methods(http.MethodPost)
	handleFunc("/login", handler.NewLoginHandler()).Methods(http.MethodPost)
	handleFunc("/create_room", handler.NewCreateRoomHandler()).Methods(http.MethodPost)
	handleFunc("/join_room", handler.NewJoinRoomHandler()).Methods(http.MethodPost)
	handleFunc("/leave_room", handler.NewLeaveRoomHandler()).Methods(http.MethodPost)
	handleFunc("/message", handler.NewMessageHandler()).Methods(http.MethodPost)

	handleFunc("/websocket", handler.NewActivateWebsocketHandler(upgrader))

	addr := fmt.Sprintf("%s:%d", "0.0.0.0", 19999)
	serv := &http.Server{
		Handler: httpHandler(router),
		Addr:    addr,
	}

	go func() {
		log.Printf("start server on %s", addr)
		if err := serv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Printf("error: %+v", err)
			}
		}
	}()

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	log.Printf("received signal: %v", <-signalCh)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Printf("shutdown: %s", addr)
	if err := serv.Shutdown(ctx); err != nil {
		panic(err)
	}

	log.Print("complete!")
}

func httpHandler(router *mux.Router) http.Handler {
	allowedHeaders := handlers.AllowedHeaders([]string{
		handler.HeaderContentType,
		handler.HeaderAccessToken,
		handler.HeaderAuthToken,
	})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	})
	const PreflightMaxAgeSec = 10 * 60 // 10 minutes
	maxAge := handlers.MaxAge(PreflightMaxAgeSec)
	return handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods, maxAge)(router)
}

func handleFunc(path string, handler http.Handler) *mux.Route {
	return router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		handler.ServeHTTP(w, r)
	})
}
