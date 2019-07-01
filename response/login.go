package response

import "encoding/json"

type Login struct {
	AccessToken string `json:"access_token"`
}

func (l *Login) Encode() ([]byte, error) {
	return json.Marshal(l)
}
