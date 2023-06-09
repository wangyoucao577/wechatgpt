package main

import (
	"flag"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/wangyoucao577/wechatgpt/wechat"
)

var wechatFlags struct {
	token          string
	usersWhitelist string
}

func init() {
	flag.StringVar(&wechatFlags.token, "wechat_token", "", "Your token to verify wechat platform requests.")
	flag.StringVar(&wechatFlags.usersWhitelist, "wechat_users_whitelist", "", "wechat users(OpenIDs) that permitted to talk to, split by ','.")
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

func replyMessageResponse(wxResp *wechat.Message, c *gin.Context) {
	wxResp.CreateTime = time.Now().Unix()

	glog.V(1).Infof("wechat response: %s\n", wxResp.String())

	if b, err := wxResp.Marshal(); err != nil {
		c.String(http.StatusBadGateway, "xml marshal failed, err %v", err)
	} else {
		c.String(http.StatusOK, string(b))
	}
}

func wxMessageHandler(c *gin.Context) {
	if gin.Mode() == gin.ReleaseMode { // validate all wechat requests on release
		token := wechatFlags.token

		signature := c.Query("signature")
		timestamp := c.Query("timestamp")
		nonce := c.Query("nonce")
		echostr := c.Query("echostr")
		if _, err := wechat.New(token).Validate(signature, timestamp, nonce, echostr); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	var wxReq, wxResp wechat.Message

	if err := wxReq.Decode(c.Request.Body); err != nil {
		c.String(http.StatusBadRequest, "decode wx message failed")
		return
	}
	glog.V(1).Infof("wechat request: %s\n", wxReq.String())

	wxResp.FromUserName.Value = wxReq.ToUserName.Value
	wxResp.ToUserName.Value = wxReq.FromUserName.Value
	wxResp.MsgType.Value = wechat.MessageTypeText // only support text response at the moment

	if wxReq.Content != nil && (wxReq.Content.Value == "help" || wxReq.Content.Value == "帮助") { // return help information
		wxResp.Content = &wechat.Content{Value: "回复\"help\"或\"帮助\"获取帮助\n回复任意内容开启对话\n问题问完过一会回复\"1\"来获取答案"}
		replyMessageResponse(&wxResp, c)
		return
	}

	if gin.Mode() == gin.ReleaseMode && len(wechatFlags.usersWhitelist) > 0 { // only talk to whitelist users on release
		usersWhitelist := wechatFlags.usersWhitelist
		if !strings.Contains(usersWhitelist, wxReq.FromUserName.Value) {
			wxResp.Content = &wechat.Content{Value: "不想跟你说话"}
			replyMessageResponse(&wxResp, c)
			return
		}
	}

	if wxReq.Content != nil && wxReq.Content.Value == "1" { // fetch answers
		wxResp.Content = &wechat.Content{Value: "没有了"} // default value

		if answersChanAny, ok := answersMap.Load(wxReq.FromUserName.Value); ok {
			answersChan := answersChanAny.(chan string)
			select {
			case answer := <-answersChan:
				wxResp.Content.Value = answer
			default:
				glog.Warningf("no answer available")
			}
		}

		replyMessageResponse(&wxResp, c)
		return
	}

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
	// wxResp.Content = &wechat.Content{}
	// wxResp.Content.Value = chatgpt(questionForGPT, time.Duration(time.Millisecond*4500)) // almost 5 seconds due to wechat's limitation
	questionsChan <- question{user: wxReq.FromUserName.Value, content: questionForGPT}

	for i := 0; i < 4; i++ { // try to fetch and answer in 4 seconds
		time.Sleep(time.Second)

		if answersChanAny, ok := answersMap.Load(wxReq.FromUserName.Value); ok {
			answersChan := answersChanAny.(chan string)
			select {
			case answer := <-answersChan:
				wxResp.Content = &wechat.Content{Value: answer}
			default:
				glog.V(1).Infof("no answer available")
			}
		}
		if wxResp.Content == nil { // no answer available
			continue
		}

		replyMessageResponse(&wxResp, c)
		return
	}

	c.String(http.StatusOK, "success") // so that wechat won't retry
}
