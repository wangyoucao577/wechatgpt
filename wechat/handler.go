// Package wechat provides wechat platform message processing.
package wechat

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"sort"
	"strings"
)

// Handler reprensents wechat handler.
type Handler struct {
	token string
}

// New creates a new wechat handler.
func New(token string) *Handler {
	return &Handler{
		token: token,
	}
}

// Validate validates whether request from wechat platform or not.
// https://developers.weixin.qq.com/doc/offiaccount/Getting_Started/Getting_Started_Guide.html
func (h Handler) Validate(signature, timestamp, nonce, echostr string) (string, error) {

	l := []string{h.token, timestamp, nonce}
	sort.Strings(l)
	hashcode := sha1.Sum([]byte(strings.Join(l, "")))
	hashcodeStr := hex.EncodeToString(hashcode[:])

	if hashcodeStr == signature {
		return echostr, nil
	}

	return "", errors.New("bad request")
}
