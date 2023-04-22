package wechat

import (
	"strings"
	"testing"
)

func TestMessageDecode(t *testing.T) {
	cases := []struct {
		m Message
		s string
	}{
		{
			m: Message{},
			s: "<xml></xml>",
		},
		{
			m: Message{
				FromUserName: FromUserName{Value: "12345"},
				ToUserName:   ToUserName{Value: "abcde"},
				MsgType:      MsgType{Value: "text"},
			},
			s: "<xml><FromUserName><![CDATA[12345]]></FromUserName><ToUserName><![CDATA[abcde]]></ToUserName><MsgType><![CDATA[text]]></MsgType></xml>",
		},
		{
			m: Message{
				FromUserName: FromUserName{Value: "12345"},
				ToUserName:   ToUserName{Value: "abcde"},
				MsgType:      MsgType{Value: "text"},
				Content:      &Content{Value: "aaabbb"},
				CreateTime:   123456789,
				MsgId:        987654321,
			},
			s: "<xml><FromUserName><![CDATA[12345]]></FromUserName><ToUserName><![CDATA[abcde]]></ToUserName><MsgType><![CDATA[text]]></MsgType><Content><![CDATA[aaabbb]]></Content><CreateTime>123456789</CreateTime><MsgId>987654321</MsgId></xml>",
		},
		{
			m: Message{
				FromUserName: FromUserName{Value: "12345"},
				ToUserName:   ToUserName{Value: "abcde"},
				MsgType:      MsgType{Value: "voice"},
				CreateTime:   123456789,
				MsgId:        987654321,
				Recognition:  &Recognition{Value: "cccddd"},
			},
			s: "<xml><FromUserName><![CDATA[12345]]></FromUserName><ToUserName><![CDATA[abcde]]></ToUserName><MsgType><![CDATA[voice]]></MsgType><Recognition><![CDATA[cccddd]]></Recognition><CreateTime>123456789</CreateTime><MsgId>987654321</MsgId></xml>",
		},
	}

	for _, c := range cases {
		var m Message
		if err := m.Decode(strings.NewReader(c.s)); err != nil {
			t.Error(err)
		}

		if m.ToUserName.Value != c.m.ToUserName.Value ||
			m.FromUserName.Value != c.m.FromUserName.Value ||
			m.MsgType.Value != c.m.MsgType.Value ||
			m.CreateTime != c.m.CreateTime ||
			m.MsgId != c.m.MsgId {
			t.Errorf("\n%+v\n != \n%+v\n", m, c.m)
		}

		if m.Content != nil && c.m.Content != nil && m.Content.Value != c.m.Content.Value {
			t.Errorf("Content \n%+v\n != \n%+v\n", m.Content, c.m.Content)
		}

		if m.Recognition != nil && c.m.Recognition != nil && m.Recognition.Value != c.m.Recognition.Value {
			t.Errorf("Recognition \n%+v\n != \n%+v\n", m.Recognition, c.m.Recognition)
		}

	}
}

func TestMessageMarshal(t *testing.T) {
	cases := []struct {
		m Message
		s string
	}{
		{
			m: Message{},
			s: "<xml><ToUserName></ToUserName><FromUserName></FromUserName><MsgType></MsgType><CreateTime>0</CreateTime></xml>",
		},
		{
			m: Message{
				FromUserName: FromUserName{Value: "12345"},
				ToUserName:   ToUserName{Value: "abcde"},
				MsgType:      MsgType{Value: "text"},
				Content:      &Content{Value: "aaabbb"},
				CreateTime:   123456789,
				MsgId:        987654321,
			},
			s: "<xml><ToUserName><![CDATA[abcde]]></ToUserName><FromUserName><![CDATA[12345]]></FromUserName><MsgType><![CDATA[text]]></MsgType><Content><![CDATA[aaabbb]]></Content><CreateTime>123456789</CreateTime><MsgId>987654321</MsgId></xml>",
		},
		{
			m: Message{
				FromUserName: FromUserName{Value: "12345"},
				ToUserName:   ToUserName{Value: "abcde"},
				MsgType:      MsgType{Value: "voice"},
				Recognition:  &Recognition{Value: "cccddd"},
				CreateTime:   123456789,
				MsgId:        987654321,
			},
			s: "<xml><ToUserName><![CDATA[abcde]]></ToUserName><FromUserName><![CDATA[12345]]></FromUserName><MsgType><![CDATA[voice]]></MsgType><CreateTime>123456789</CreateTime><MsgId>987654321</MsgId><Recognition><![CDATA[cccddd]]></Recognition></xml>",
		},
	}

	for _, c := range cases {
		var s string
		if b, err := c.m.Marshal(); err != nil {
			t.Error(err)
		} else {
			s = string(b)
		}

		if s != c.s {
			t.Errorf("\n%s\n != \n%s\n", s, c.s)
		}
	}
}
