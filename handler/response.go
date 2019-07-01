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
	r.response(200, w)
}

func (r *Response) NoContent(w http.ResponseWriter) {
	w.WriteHeader(204)
}

func (r *Response) BadRequest(w http.ResponseWriter) {
	r.response(400, w)
}

func (r *Response) Unauthorized(w http.ResponseWriter) {
	r.response(401, w)
}

func (r *Response) Forbidden(w http.ResponseWriter) {
	r.response(403, w)
}

func (r *Response) NotFound(w http.ResponseWriter) {
	r.response(404, w)
}

func (r *Response) InternalError(w http.ResponseWriter) {
	r.response(500, w)
}

func (r *Response) Unavailable(w http.ResponseWriter) {
	r.response(503, w)
}
