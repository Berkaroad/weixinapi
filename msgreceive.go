package weixinapi

import (
	"encoding/xml"
	"log"
	"time"
)

// 消息基类
type baseMessage struct {
	ToUserName   CDATAText     `xml:"ToUserName"`
	FromUserName CDATAText     `xml:"FromUserName"`
	CreateTime   time.Duration `xml:"CreateTime"`
	MsgType      CDATAText     `xml:"MsgType"`
	MsgId        int           `xml:"MsgId"`
}

// 文本响应消息
type TextResponseMessage struct {
	XMLName xml.Name `xml:"xml"`
	baseMessage
	Content CDATAText `xml:"Content"`
}

func NewTextResponseMessage(toUserName string, fromUserName string, createTime time.Duration, content string, msgId int) TextResponseMessage {
	return TextResponseMessage{baseMessage: baseMessage{ToUserName: NewCDATAText(toUserName), FromUserName: NewCDATAText(fromUserName), CreateTime: createTime, MsgType: NewCDATAText("text"), MsgId: msgId}, Content: NewCDATAText(content)}
}

// 位置响应消息
type LocationResponseMessage struct {
	XMLName xml.Name `xml:"xml"`
	baseMessage
	Location_X float64   `xml:"Location_X"`
	Location_Y float64   `xml:"Location_Y"`
	Scale      int       `xml:"Scale"`
	Label      CDATAText `xml:"Label"`
}

func NewLocationResponseMessage(toUserName string, fromUserName string, createTime time.Duration, location_X float64, location_Y float64, scale int, label string, msgId int) LocationResponseMessage {
	return LocationResponseMessage{baseMessage: baseMessage{ToUserName: NewCDATAText(toUserName), FromUserName: NewCDATAText(fromUserName), CreateTime: createTime, MsgType: NewCDATAText("location"), MsgId: msgId}, Location_X: location_X, Location_Y: location_Y, Scale: scale, Label: NewCDATAText(label)}
}

// 链接响应消息
type LinkResponseMessage struct {
	XMLName xml.Name `xml:"xml"`
	baseMessage
	Title       CDATAText `xml:"Title"`
	Description CDATAText `xml:"Description"`
	Url         CDATAText `xml:"Url"`
}

func NewLinkResponseMessage(toUserName string, fromUserName string, createTime time.Duration, title string, description string, url string, msgId int) LinkResponseMessage {
	return LinkResponseMessage{baseMessage: baseMessage{ToUserName: NewCDATAText(toUserName), FromUserName: NewCDATAText(fromUserName), CreateTime: createTime, MsgType: NewCDATAText("link"), MsgId: msgId}, Title: NewCDATAText(title), Description: NewCDATAText(description), Url: NewCDATAText(url)}
}

// Encrypt Text Message
func EncryptTextMsg(config WXAPIConfig, replyMsg TextResponseMessage, timestamp string, nonce string) (string, error) {
	bReplyMsg, _ := xml.Marshal(&replyMsg)
	sReplyMsg := string(bReplyMsg)
	log.Println("raw_message=" + sReplyMsg)
	result, err := EncryptMsg(config, sReplyMsg, timestamp, nonce)
	log.Println("encrypt_message=" + result)
	return result, err
}

// Encrypt Location Message
func EncryptLocationMsg(config WXAPIConfig, replyMsg LocationResponseMessage, timestamp string, nonce string) (string, error) {
	bReplyMsg, _ := xml.Marshal(&replyMsg)
	sReplyMsg := string(bReplyMsg)
	log.Println("raw_message=" + sReplyMsg)
	result, err := EncryptMsg(config, sReplyMsg, timestamp, nonce)
	log.Println("encrypt_message=" + result)
	return result, err
}

// Encrypt Link Message
func EncryptLinkMsg(config WXAPIConfig, replyMsg LinkResponseMessage, timestamp string, nonce string) (string, error) {
	bReplyMsg, _ := xml.Marshal(&replyMsg)
	sReplyMsg := string(bReplyMsg)
	log.Println("raw_message=" + sReplyMsg)
	result, err := EncryptMsg(config, sReplyMsg, timestamp, nonce)
	log.Println("encrypt_message=" + result)
	return result, err
}
