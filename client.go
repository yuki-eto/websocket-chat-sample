package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"websocket-chat-sample/client"

	"github.com/peterh/liner"
)

func main() {
	var (
		name       = flag.String("name", "test_name", "user name")
		loginToken = flag.String("token", "", "login token")
		roomID     = flag.String("room_id", "test_room_id", "join room_id")
	)
	flag.Parse()

	cli := client.NewClient()

	if *loginToken == "" {
		if err := cli.CreateUser(*name); err != nil {
			panic(err)
		}
	} else {
		cli.LoginToken = *loginToken
	}

	joinedRoomID, err := cli.Login()
	if err != nil {
		panic(err)
	}
	if joinedRoomID != "" {
		if err := cli.LeaveRoom(); err != nil {
			panic(err)
		}
	}

	if err := cli.HandleWebsocket(); err != nil {
		panic(err)
	}

	if err := cli.CreateRoom(*roomID); err != nil {
		panic(err)
	}

	if err := cli.JoinRoom(*roomID); err != nil {
		panic(err)
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		defer func() {
			signalCh <- syscall.SIGQUIT
		}()

		line := liner.NewLiner()
		for {
			text, err := line.Prompt(">> ")
			if err != nil {
				if err == liner.ErrPromptAborted {
					log.Print("aborted.")
				} else {
					log.Printf("error: %+v", err)
				}
				return
			}

			if text == "" {
				continue
			} else if strings.HasPrefix(text, "/q") {
				return
			}

			if err := cli.Message(text); err != nil {
				log.Printf("error: %+v", err)
				return
			}
		}
	}()

	<-signalCh

	if err := cli.LeaveRoom(); err != nil {
		panic(err)
	}

	log.Printf("leave room: %s", *roomID)
}
