package cookie

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"html"
	"net/http"
	"strconv"
	"time"
)

type CookiePacker struct {
	privateKey ed25519.PrivateKey
}

func NewCookiePacker(privateKey []byte) (*CookiePacker, error) {
	block, _ := pem.Decode(privateKey)
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &CookiePacker{privateKey: key.(ed25519.PrivateKey)}, nil
}

func (c CookiePacker) PackAndSign(payload string) http.Cookie {
	cookie := http.Cookie{}
	now := time.Now()

	cookie.Name = "token"
	cookie.HttpOnly = true

	payloadEncoded := html.EscapeString(payload)
	dateEncoded := html.EscapeString(strconv.FormatInt(now.Unix(), 10))
	signedMessage := payloadEncoded + "#" + dateEncoded
	signature := ed25519.Sign(c.privateKey, []byte(signedMessage))
	signedMessage += "#" + base64.URLEncoding.EncodeToString(signature)
	cookie.Value = signedMessage

	return cookie
}
