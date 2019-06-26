package entity

type User struct {
	ID     uint64 `json:"id"`
	Token  string `json:"token"`
	Name   string `json:"name"`
	RoomID string `json:"room_id"`
}
