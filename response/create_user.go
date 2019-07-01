package response

import "encoding/json"

type CreateUser struct {
	LoginToken string `json:"login_token"`
}

func (c *CreateUser) Encode() ([]byte, error) {
	return json.Marshal(c)
}
