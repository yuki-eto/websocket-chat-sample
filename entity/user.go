package entity

type User struct {
	ID          uint64 `json:"id"`
	LoginToken  string `json:"-"`
	Name        string `json:"name"`
	AccessToken string `json:"-"`
}
