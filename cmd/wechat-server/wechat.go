package main

import (
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wangyoucao577/wechatgpt/wechat"
)

var wechatFlags struct {
	token string
}

func init() {
	flag.StringVar(&wechatFlags.token, "token", "", "Your token to verify wechat platform requests.")
}

func wxValidationHandler(c *gin.Context) {

	token := wechatFlags.token

	// https://developers.weixin.qq.com/doc/offiaccount/Getting_Started/Getting_Started_Guide.html

	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")

	l := []string{token, timestamp, nonce}
	sort.Strings(l)
	hashcode := sha1.Sum([]byte(strings.Join(l, "")))
	hashcodeStr := hex.EncodeToString(hashcode[:])

	if hashcodeStr == signature {
		c.String(http.StatusOK, echostr)
		return
	}
	c.String(http.StatusBadRequest, "")
}

func wxMessageHandler(c *gin.Context) {
	var wxReq, wxResp wechat.Message

	if err := wxReq.Decode(c.Request.Body); err != nil {
		c.String(http.StatusBadRequest, "decode wx message failed")
		return
	}

	wxResp.FromUserName.Value = wxReq.ToUserName.Value
	wxResp.ToUserName.Value = wxReq.FromUserName.Value
	wxResp.CreateTime = time.Now().Unix()
	wxResp.MsgType.Value = wechat.MessageTypeText

	var questionForGPT string
	switch wxReq.MsgType.Value {
	case wechat.MessageTypeText:
		questionForGPT = wxReq.Content.Value
	case wechat.MessageTypeVoice:
		questionForGPT = wxReq.Recognition.Value
	default: // only process text/voice recoginzation at the moment
		c.String(http.StatusBadRequest, "unsupport message type "+wxReq.MsgType.Value)
		return
	}

	// generate response via chatgpt
	wxResp.Content.Value = chatgpt(questionForGPT, time.Duration(time.Millisecond*4900)) // almost 5 seconds due to wechat's limitation

	if b, err := wxResp.Marshal(); err != nil {
		c.String(http.StatusBadGateway, "xml marshal failed, err %v", err)
	} else {
		c.String(http.StatusOK, string(b))
	}
}
