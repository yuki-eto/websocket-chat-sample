package handler

import (
	"log"
	"net/http"
	"websocket-chat-sample/repository"
	"websocket-chat-sample/response"

	"github.com/juju/errors"
)

const (
	HeaderContentType = "content-type"
	HeaderAuthToken   = "x-authenticate-token"
	HeaderAccessToken = "x-access-token"
	FormLoginToken    = "login_token"
	FormAccessToken   = "access_token"
)

func GetUser(w http.ResponseWriter, r *http.Request, userRepo repository.UserRepository) *repository.UserInstance {
	loginToken := r.Header.Get(HeaderAuthToken)
	if loginToken == "" {
		loginToken = r.FormValue(FormLoginToken)
	}
	accessToken := r.Header.Get(HeaderAccessToken)
	if accessToken == "" {
		accessToken = r.FormValue(FormAccessToken)
	}

	log.Printf("login: %s, access: %s", loginToken, accessToken)

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
