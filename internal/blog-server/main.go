package blogserver

import (
	"net/http"

	"example.org/nn/kaftinker/internal/blog-server/handler"
)

func Serve(addr string, h handler.Handler) error {
	return http.ListenAndServe(addr, nil)
}
