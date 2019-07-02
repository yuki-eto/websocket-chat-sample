package handler

import (
	"net/http"
	"websocket-chat-sample/repository"
	"websocket-chat-sample/response"

	"github.com/juju/errors"
)

const (
	HeaderAuthToken   = "x-authenticate-token"
	HeaderAccessToken = "x-access-token"
)

func GetUser(w http.ResponseWriter, r *http.Request, userRepo repository.UserRepository) *repository.UserInstance {
	loginToken := r.Header.Get(HeaderAuthToken)
	accessToken := r.Header.Get(HeaderAccessToken)

	user, err := userRepo.FindByToken(loginToken)
	if err != nil {
		res := &Response{
			&response.ErrorResponse{Error: err},
		}
		if errors.IsNotFound(err) {
			res.NotFound(w)
		} else {
			res.InternalError(w)
		}
		return nil
	}

	if user.AccessToken != accessToken {
		res := &Response{
			&response.ErrorResponse{Error: errors.NewUnauthorized(nil, "access token is not match")},
		}
		res.Unauthorized(w)
		return nil
	}

	return user
}
