package response

import (
	"encoding/json"
	"fmt"
)

type ErrorResponse struct {
	Error error `json:"error"`
}

func (e *ErrorResponse) Encode() ([]byte, error) {
	dict := map[string]interface{}{
		"error_trace": fmt.Sprintf("%+v", e.Error),
	}
	return json.Marshal(dict)
}
