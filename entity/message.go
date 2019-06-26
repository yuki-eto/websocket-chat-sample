package entity

import "time"

type Message struct {
	RoomID string     `json:"-"`
	UserID uint64     `json:"user_id"`
	Text   string     `json:"text"`
	Time   *time.Time `json:"time"`
}
