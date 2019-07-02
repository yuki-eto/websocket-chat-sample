package handler

import (
	"log"
	"net/http"
	"websocket-chat-sample/response"
)

type Response struct {
	body response.Body
}

func (r *Response) response(code int, w http.ResponseWriter) {
	w.WriteHeader(code)

	b, err := r.body.Encode()
	if err != nil {
		log.Printf("error: %+v", err)
		return
	}
	if _, err := w.Write(b); err != nil {
		log.Printf("error: %+v", err)
	}
}

func (r *Response) Ok(w http.ResponseWriter) {
	r.response(http.StatusOK, w)
}

func (r *Response) NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func (r *Response) BadRequest(w http.ResponseWriter) {
	r.response(http.StatusBadRequest, w)
}

func (r *Response) Unauthorized(w http.ResponseWriter) {
	r.response(http.StatusUnauthorized, w)
}

func (r *Response) Forbidden(w http.ResponseWriter) {
	r.response(http.StatusForbidden, w)
}

func (r *Response) NotFound(w http.ResponseWriter) {
	r.response(http.StatusNotFound, w)
}

func (r *Response) Conflict(w http.ResponseWriter) {
	r.response(http.StatusConflict, w)
}

func (r *Response) InternalError(w http.ResponseWriter) {
	r.response(http.StatusInternalServerError, w)
}

func (r *Response) Unavailable(w http.ResponseWriter) {
	r.response(http.StatusServiceUnavailable, w)
}
