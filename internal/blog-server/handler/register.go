package handler

import (
	"crypto/sha512"
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"example.org/nn/kaftinker/internal/blog-server/utils"
	"example.org/nn/kaftinker/internal/types"
)

type RegisterPageData struct {
	UsernameTaken       bool
	UsernameNotProvided bool
	PasswordNotProvided bool
}

func (h Handler) getRegister(w http.ResponseWriter, r *http.Request) {
	if _, err := utils.GetAuthenticatedUser(h.CookieUnpacker, r); err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	if err := h.Template.ExecuteTemplate(w, "register-page", RegisterPageData{}); err != nil {
		log.Println(err)
	}
}

func (h Handler) postRegister(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" {
		if err := h.Template.ExecuteTemplate(w, "register-page", RegisterPageData{UsernameNotProvided: true}); err != nil {
			log.Println(err)
		}
		return
	}
	if password == "" {
		if err := h.Template.ExecuteTemplate(w, "register-page", RegisterPageData{PasswordNotProvided: true}); err != nil {
			log.Println(err)
		}
		return
	}

	_, err := h.Storage.GetUserByUsername(h.Ctx, username)
	if err != sql.ErrNoRows {
		if err := h.Template.ExecuteTemplate(w, "register-page", RegisterPageData{UsernameTaken: true}); err != nil {
			log.Println(err)
		}
		return
	}

	hashedPassword := sha512.Sum512([]byte(password + h.PasswordSalt))

	newUser := types.User{
		Username: username,
		Password: hashedPassword[:],
	}
	userId, err := h.Storage.CreateUser(h.Ctx, newUser)
	if err != nil {
		log.Println(err)

		utils.WriteInternalError(w)
		return
	}

	payload := strconv.FormatInt(userId, 10)

	cookie := h.CookiePacker.PackAndSign(payload)
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}
