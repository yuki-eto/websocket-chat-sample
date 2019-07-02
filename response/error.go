package response

import (
	"encoding/json"
	"fmt"
)

var isDebug = true

const (
	ErrorCodeUnknown          = 100000001
	ErrorCodeAlreadyExistRoom = 100001001
)

type ErrorResponse struct {
	Error error `json:"error"`
	Code  int   `json:"code"`
}

func (e *ErrorResponse) Encode() ([]byte, error) {
	if e.Code == 0 {
		e.Code = ErrorCodeUnknown
	}

	dict := map[string]interface{}{
		"error": e.Error.Error(),
		"code":  e.Code,
	}
	if isDebug == true {
		dict["error"] = fmt.Sprintf("%+v", e.Error)
	}

	return json.Marshal(dict)
}
