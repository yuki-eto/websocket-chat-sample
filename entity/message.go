package entity

import "time"

type Message struct {
	RoomID string     `json:"-"`
	Name   string     `json:"name"`
	Text   string     `json:"text"`
	Time   *time.Time `json:"time"`
}
