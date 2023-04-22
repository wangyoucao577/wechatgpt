package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
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

	if s, err := wechat.New(token).Validate(signature, timestamp, nonce, echostr); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	} else {
		c.String(http.StatusOK, s)
	}
}

func wxMessageHandler(c *gin.Context) {
	var wxReq, wxResp wechat.Message

	if err := wxReq.Decode(c.Request.Body); err != nil {
		c.String(http.StatusBadRequest, "decode wx message failed")
		return
	}
	glog.V(1).Infof("wechat request: %s\n", wxReq.String())

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
	wxResp.Content = &wechat.Content{}
	wxResp.Content.Value = chatgpt(questionForGPT, time.Duration(time.Millisecond*4500)) // almost 5 seconds due to wechat's limitation

	glog.V(1).Infof("wechat response: %s\n", wxResp.String())

	if b, err := wxResp.Marshal(); err != nil {
		c.String(http.StatusBadGateway, "xml marshal failed, err %v", err)
	} else {
		c.String(http.StatusOK, string(b))
	}
}
