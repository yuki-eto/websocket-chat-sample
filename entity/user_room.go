package entity

import "time"

type UserRoom struct {
	UserID   uint64     `json:"user_id"`
	RoomID   string     `json:"room_id"`
	JoinedAt *time.Time `json:"joined_at"`
}
