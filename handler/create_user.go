package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"websocket-chat-sample/entity"
	"websocket-chat-sample/repository"
	"websocket-chat-sample/request"
	"websocket-chat-sample/response"

	"github.com/google/uuid"
	"github.com/juju/errors"
)

type CreateUser struct {
	user repository.UserRepository
}

func NewCreateUserHandler() http.Handler {
	return &CreateUser{
		user: repository.NewUserRepository(),
	}
}

func (c *CreateUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req request.CreateUser
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		res := &Response{
			&response.ErrorResponse{Error: err},
		}
		res.BadRequest(w)
		return
	}
	if req.Name == "" {
		res := &Response{
			&response.ErrorResponse{Error: errors.NewBadRequest(nil, "bad request")},
		}
		res.BadRequest(w)
		return
	}

	u, err := uuid.NewRandom()
	if err != nil {
		res := &Response{
			&response.ErrorResponse{Error: err},
		}
		res.InternalError(w)
		return
	}

	loginToken := u.String()
	user := repository.NewUserInstance(&entity.User{
		Name:       req.Name,
		LoginToken: loginToken,
	})
	if err := c.user.Save(user); err != nil {
		res := &Response{
			&response.ErrorResponse{Error: err},
		}
		res.InternalError(w)
		return
	}

	res := &Response{
		&response.CreateUser{LoginToken: loginToken},
	}
	res.Ok(w)

	log.Printf("created user : %+v", user.User)
}
