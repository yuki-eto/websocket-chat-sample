package response

import (
	"websocket-chat-sample/repository"
)

type StreamType string

const (
	StreamTypeChat  StreamType = "chat"
	StreamTypeJoin  StreamType = "join"
	StreamTypeLeave StreamType = "leave"
)

type Stream struct {
	Type StreamType `json:"type"`

	Message *repository.MessageInstance `json:"message"`
	User    *repository.UserInstance    `json:"user"`
}
