package wechat

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

// Message Types
const (
	MessageTypeText  = "text"
	MessageTypeVoice = "voice"
)

// ToUserName represents ToUserName structure of wechat message.
type ToUserName struct {
	XMLName xml.Name `xml:"ToUserName"`
	Value   string   `xml:",cdata"`
}

// FromUserName represents FromUserName structure of wechat message.
type FromUserName struct {
	XMLName xml.Name `xml:"FromUserName"`
	Value   string   `xml:",cdata"`
}

// MsgType represents MsgType structure of wechat message.
type MsgType struct {
	XMLName xml.Name `xml:"MsgType"`
	Value   string   `xml:",cdata"`
}

// Content represents Content structure of wechat message.
type Content struct {
	XMLName xml.Name `xml:"Content,omitempty"`
	Value   string   `xml:",cdata"`
}

// Recognition represents Recognition structure of wechat message.
type Recognition struct {
	XMLName xml.Name `xml:"Recognition,omitempty"`
	Value   string   `xml:",cdata"`
}

// Message represents wechat platform message structure.
type Message struct {
	XMLName      xml.Name     `xml:"xml"`
	ToUserName   ToUserName   `xml:"ToUserName"`
	FromUserName FromUserName `xml:"FromUserName"`
	MsgType      MsgType      `xml:"MsgType"`
	Content      *Content     `xml:"Content,omitempty"`

	CreateTime int64 `xml:"CreateTime"`
	MsgId      int64 `xml:"MsgId,omitempty"`

	Recognition *Recognition `xml:"Recognition,omitempty"`
}

// Marshal marshals wechat message.
func (m *Message) Marshal() ([]byte, error) {
	return xml.Marshal(*m)
}

// Unmarshal unmarshals wechat message.
func (m *Message) Unmarshal(b []byte) error {
	return m.Decode(bytes.NewReader(b))
}

// Decode decodes wechat message.
func (m *Message) Decode(r io.Reader) error {
	return xml.NewDecoder(r).Decode(m)
}

// String generates string of message.
func (m Message) String() string {
	s := fmt.Sprintf("ToUserName %s FromUserName %s MsgType %s CreateTime %d MsgId %d", m.ToUserName.Value, m.FromUserName.Value, m.MsgType.Value, m.CreateTime, m.MsgId)
	if m.Content != nil {
		s += fmt.Sprintf(" Content %s", m.Content.Value)
	}
	if m.Recognition != nil {
		s += fmt.Sprintf(" Recognition %s", m.Recognition.Value)
	}
	return s
}
