package utils

import (
	"io"
	"net/http"
	"strconv"

	"example.org/nn/kaftinker/internal/blog-server/utils/cookie"
)

func WriteInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)

	io.WriteString(w, "Произошла неизвестная ошибка")
}

func GetAuthenticatedUser(cookieUnpacker *cookie.CookieUnpacker, r *http.Request) (int64, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return 0, err
	}

	idMarshaled, err := cookieUnpacker.VerifyAndUnpack(*cookie)
	if err != nil {
		return 0, err
	}
	// idMarshaled = bytes.TrimRight(idMarshaled, string(rune(0)))

	id, err := strconv.ParseInt(string(idMarshaled), 10, 64)

	return id, err
}
