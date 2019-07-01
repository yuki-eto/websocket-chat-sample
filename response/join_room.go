package response

import (
	"encoding/json"
	"websocket-chat-sample/repository"
)

type JoinRoom struct {
	Room     *repository.RoomInstance     `json:"room"`
	Users    []*repository.UserInstance   `json:"users"`
	Messages *repository.MessagesInstance `json:"messages"`
}

func (j *JoinRoom) Encode() ([]byte, error) {
	return json.Marshal(j)
}
