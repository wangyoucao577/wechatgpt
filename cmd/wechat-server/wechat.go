package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"flag"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

const (
	wxMessageTypeText  = "text"
	wxMessageTypeVoice = "voice"
)

type wxMessage struct {
	XMLName    xml.Name `xml:"xml"`
	ToUserName struct {
		XMLName xml.Name `xml:"ToUserName"`
		Value   string   `xml:",cdata"`
	}
	FromUserName struct {
		XMLName xml.Name `xml:"FromUserName"`
		Value   string   `xml:",cdata"`
	}
	MsgType struct {
		XMLName xml.Name `xml:"MsgType"`
		Value   string   `xml:",cdata"`
	}
	Content struct {
		XMLName xml.Name `xml:"Content"`
		Value   string   `xml:",cdata"`
	}

	CreateTime int64 `xml:"CreateTime"`
	MsgId      int64 `xml:"MsgId,omitempty"`

	// Recognition struct {
	// 	XMLName  xml.Name `xml:"Recognition,omitempty"`
	// 	Value string   `xml:",chardata,omitempty"`
	// }
}

func wxMessageHandler(c *gin.Context) {
	var wxReq, wxResp wxMessage

	if err := xml.NewDecoder(c.Request.Body).Decode(&wxReq); err != nil {
		c.String(http.StatusBadRequest, "decode wx message failed")
		return
	}

	if wxReq.MsgType.Value != wxMessageTypeText {
		c.String(http.StatusBadRequest, "unsupport message type "+wxReq.MsgType.Value)
		return
	}

	wxResp.FromUserName.Value = wxReq.ToUserName.Value
	wxResp.ToUserName.Value = wxReq.FromUserName.Value
	wxResp.CreateTime = time.Now().Unix()
	wxResp.MsgType.Value = wxMessageTypeText

	// TODO: generate response
	wxResp.Content.Value = wxReq.Content.Value

	if b, err := xml.Marshal(wxResp); err != nil {
		c.String(http.StatusBadGateway, "xml marshal failed, err %v", err)
	} else {
		c.String(http.StatusOK, string(b))
	}
}
