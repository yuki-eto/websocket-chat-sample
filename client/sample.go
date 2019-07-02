package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"websocket-chat-sample/handler"
	"websocket-chat-sample/request"
	"websocket-chat-sample/response"

	"github.com/gorilla/websocket"
	"github.com/juju/errors"
)

type Client struct {
	LoginToken   string
	AccessToken  string
	joinedRoomID string
	conn         *websocket.Conn
	initialState string

	streams []string
}

func NewClient() *Client {
	return &Client{
		streams: []string{},
	}
}

func (c *Client) CreateUser(name string) error {
	req := &request.CreateUser{
		Name: name,
	}
	b, err := json.Marshal(req)
	if err != nil {
		return errors.Trace(err)
	}
	body := bytes.NewBuffer(b)
	resp, err := http.Post("http://localhost:19999/create_user", "application/json", body)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode >= 400 {
		b, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(b))
	}

	var res response.CreateUser
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&res); err != nil {
		return errors.Trace(err)
	}

	c.LoginToken = res.LoginToken
	return nil

}

func (c *Client) Login() (string, error) {
	req, err := http.NewRequest("POST", "http://localhost:19999/login", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set(handler.HeaderAuthToken, c.LoginToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Trace(err)
	}
	if resp.StatusCode >= 400 {
		b, _ := ioutil.ReadAll(resp.Body)
		return "", errors.New(string(b))
	}

	var res response.Login
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&res); err != nil {
		return "", errors.Trace(err)
	}

	c.AccessToken = res.AccessToken
	return res.JoinedRoomID, nil
}

func (c *Client) HandleWebsocket() error {
	u := url.URL{Scheme: "ws", Host: "localhost:19999", Path: "/websocket"}
	log.Printf("connecting to %s", u.String())

	header := http.Header{}
	header.Set(handler.HeaderAuthToken, c.LoginToken)
	header.Set(handler.HeaderAccessToken, c.AccessToken)
	conn, resp, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		return errors.Trace(err)
	}
	c.conn = conn

	if resp.StatusCode >= 400 {
		b, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(b))
	}

	go func() {
		defer func() {
			if err := conn.Close(); err != nil {
				log.Printf("close error: %+v", err)
			}
		}()

		i := 0
		for {
			i++
			log.Printf("loop: %d", i)

			_, b, err := conn.ReadMessage()
			if err != nil {
				log.Printf("error: %+v", err)
				return
			}
			log.Printf(string(b))
			c.streams = append(c.streams, string(b))
		}
	}()

	return nil
}

func (c *Client) CreateRoom(roomID string) error {
	buf := bytes.NewBuffer([]byte(""))
	encoder := json.NewEncoder(buf)

	createRoom := &request.CreateRoom{RoomID: roomID}
	if err := encoder.Encode(createRoom); err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", "http://localhost:19999/create_room", buf)
	if err != nil {
		return errors.Trace(err)
	}
	req.Header.Set(handler.HeaderAuthToken, c.LoginToken)
	req.Header.Set(handler.HeaderAccessToken, c.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Trace(err)
	}

	if resp.StatusCode >= 400 {
		if resp.StatusCode != http.StatusConflict {
			return nil
		}
		b, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(b))
	}

	return nil
}

func (c *Client) JoinRoom(roomID string) error {
	buf := bytes.NewBuffer([]byte(""))
	encoder := json.NewEncoder(buf)

	joinRoom := &request.JoinRoom{RoomID: roomID}
	if err := encoder.Encode(joinRoom); err != nil {
		return errors.Trace(err)
	}
	req, err := http.NewRequest("POST", "http://localhost:19999/join_room", buf)
	if err != nil {
		return errors.Trace(err)
	}
	req.Header.Set(handler.HeaderAuthToken, c.LoginToken)
	req.Header.Set(handler.HeaderAccessToken, c.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Trace(err)
	}
	if resp.StatusCode >= 400 {
		b, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(b))
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Trace(err)
	}

	c.initialState = string(b)
	return nil
}

func (c *Client) Message(text string) error {
	buf := bytes.NewBuffer([]byte(""))
	encoder := json.NewEncoder(buf)

	msg := &request.Message{Text: text}
	if err := encoder.Encode(msg); err != nil {
		return errors.Trace(err)
	}
	req, err := http.NewRequest("POST", "http://localhost:19999/message", buf)
	if err != nil {
		return errors.Trace(err)
	}
	req.Header.Set(handler.HeaderAuthToken, c.LoginToken)
	req.Header.Set(handler.HeaderAccessToken, c.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Trace(err)
	}
	if resp.StatusCode >= 400 {
		b, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(b))
	}

	return nil
}

func (c *Client) LeaveRoom() error {
	req, err := http.NewRequest("POST", "http://localhost:19999/leave_room", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set(handler.HeaderAuthToken, c.LoginToken)
	req.Header.Set(handler.HeaderAccessToken, c.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Trace(err)
	}
	if resp.StatusCode >= 400 {
		b, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(b))
	}
	return nil
}
