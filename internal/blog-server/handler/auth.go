package handler

import (
	"crypto/sha512"
	"crypto/subtle"
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"example.org/nn/kaftinker/internal/blog-server/utils"
)

var authPageTemplate string

type AuthPageData struct {
	UserNotExists       bool
	WrongPassword       bool
	UsernameNotProvided bool
	PasswordNotProvided bool
}

func (h Handler) getAuth(w http.ResponseWriter, r *http.Request) {
	if _, err := utils.GetAuthenticatedUser(h.CookieUnpacker, r); err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	if err := h.Template.ExecuteTemplate(w, "auth-page", AuthPageData{}); err != nil {
		log.Println(err)
	}
}

func (h Handler) postAuth(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" {
		if err := h.Template.ExecuteTemplate(w, "auth-page", AuthPageData{UsernameNotProvided: true}); err != nil {
			log.Println(err)
		}
		return
	}
	if password == "" {
		if err := h.Template.ExecuteTemplate(w, "auth-page", AuthPageData{PasswordNotProvided: true}); err != nil {
			log.Println(err)
		}
		return
	}

	user, err := h.Storage.GetUserByUsername(h.Ctx, username)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
		}
		if err := h.Template.ExecuteTemplate(w, "auth-page", AuthPageData{UserNotExists: true}); err != nil {
			log.Println(err)
		}
		return
	}

	hashedPassword := sha512.Sum512([]byte(password + h.PasswordSalt))
	if subtle.ConstantTimeCompare(user.Password[:], hashedPassword[:]) != 1 {
		if err := h.Template.ExecuteTemplate(w, "auth-page", AuthPageData{WrongPassword: true}); err != nil {
			log.Println(err)
		}
		return
	}

	payload := strconv.FormatInt(user.Id, 10)

	cookie := h.CookiePacker.PackAndSign(payload)
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}
