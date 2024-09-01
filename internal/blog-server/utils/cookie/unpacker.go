package cookie

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"html"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	ErrTokensNotEnough = errors.New("Bad tokens count")
	ErrMalformedCookie = errors.New("Malformed cookie")
	ErrCookieExpired   = errors.New("Cookie expired")
)

type CookieUnpacker struct {
	publicKey ed25519.PublicKey
}

func NewCookieUnpacker(publicKey []byte) (*CookieUnpacker, error) {
	block, _ := pem.Decode(publicKey)
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &CookieUnpacker{publicKey: key.(ed25519.PublicKey)}, nil
}

func (v CookieUnpacker) VerifyAndUnpack(cookie http.Cookie) (string, error) {
	toks := strings.Split(cookie.Value, "#")
	if len(toks) < 3 {
		return "", ErrTokensNotEnough
	}

	signature, err := base64.URLEncoding.DecodeString(toks[2])
	if err != nil {
		return "", err
	}
	signedMessage := toks[0] + "#" + toks[1]

	if !ed25519.Verify(v.publicKey, []byte(signedMessage), []byte(signature)) {
		return "", ErrMalformedCookie
	}

	payload := html.UnescapeString(toks[0])
	unixTimeBytes := html.UnescapeString(toks[1])
	unixTimeNumber, err := strconv.ParseInt(string(unixTimeBytes), 10, 64)
	if err != nil {
		return "", err
	}
	unixTime := time.Unix(unixTimeNumber, 0)

	if unixTime.Add(COOKIE_LONGEVITY).Before(time.Now()) {
		return "", ErrCookieExpired
	}

	return payload, nil
}
