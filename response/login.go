package response

import "encoding/json"

type Login struct {
	AccessToken  string `json:"access_token"`
	JoinedRoomID string `json:"joined_room_id"`
}

func (l *Login) Encode() ([]byte, error) {
	return json.Marshal(l)
}
