package handler

import (
	"log"
	"net/http"
	"websocket-chat-sample/repository"
	"websocket-chat-sample/response"

	"github.com/google/uuid"
	"github.com/juju/errors"
)

type Login struct {
	user repository.UserRepository
}

func NewLoginHandler() http.Handler {
	return &Login{
		user: repository.NewUserRepository(),
	}
}

func (h *Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get(HeaderAuthToken)
	log.Printf("login_token: %s", token)

	user, err := h.user.FindByToken(token)
	if err != nil {
		res := &Response{
			&response.ErrorResponse{Error: err},
		}
		if errors.IsNotFound(err) {
			res.NotFound(w)
		} else {
			res.InternalError(w)
		}
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

	accessToken := u.String()
	user.AccessToken = accessToken
	if err := h.user.Save(user); err != nil {
		res := &Response{
			&response.ErrorResponse{Error: err},
		}
		res.InternalError(w)
		return
	}

	res := &Response{
		&response.Login{AccessToken: accessToken},
	}
	res.Ok(w)
}
