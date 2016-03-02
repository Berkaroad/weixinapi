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
}

// 文本响应消息
type TextResponseMessage struct {
	XMLName xml.Name `xml:"xml"`
	baseMessage
	MsgId   int       `xml:"MsgId"`
	Content CDATAText `xml:"Content"`
}

func NewTextResponseMessage(toUserName string, fromUserName string, createTime time.Duration, content string, msgId int) TextResponseMessage {
	return TextResponseMessage{baseMessage: baseMessage{ToUserName: NewCDATAText(toUserName), FromUserName: NewCDATAText(fromUserName), CreateTime: createTime, MsgType: NewCDATAText("text")}, Content: NewCDATAText(content), MsgId: msgId}
}

func EncryptTextMsg(config WXAPIConfig, replyMsg TextResponseMessage, timestamp string, nonce string) (string, error) {
	bReplyMsg, _ := xml.Marshal(&replyMsg)
	sReplyMsg := string(bReplyMsg)
	log.Println("raw_message=" + sReplyMsg)
	result, err := EncryptMsg(config, sReplyMsg, timestamp, nonce)
	log.Println("encrypt_message=" + result)
	return result, err
}

// 位置响应消息
type LocationResponseMessage struct {
	XMLName xml.Name `xml:"xml"`
	baseMessage
	MsgId      int       `xml:"MsgId"`
	Location_X float64   `xml:"Location_X"`
	Location_Y float64   `xml:"Location_Y"`
	Scale      int       `xml:"Scale"`
	Label      CDATAText `xml:"Label"`
}

func NewLocationResponseMessage(toUserName string, fromUserName string, createTime time.Duration, location_X float64, location_Y float64, scale int, label string, msgId int) LocationResponseMessage {
	return LocationResponseMessage{baseMessage: baseMessage{ToUserName: NewCDATAText(toUserName), FromUserName: NewCDATAText(fromUserName), CreateTime: createTime, MsgType: NewCDATAText("location")}, MsgId: msgId, Location_X: location_X, Location_Y: location_Y, Scale: scale, Label: NewCDATAText(label)}
}

func EncryptLocationMsg(config WXAPIConfig, replyMsg LocationResponseMessage, timestamp string, nonce string) (string, error) {
	bReplyMsg, _ := xml.Marshal(&replyMsg)
	sReplyMsg := string(bReplyMsg)
	log.Println("raw_message=" + sReplyMsg)
	result, err := EncryptMsg(config, sReplyMsg, timestamp, nonce)
	log.Println("encrypt_message=" + result)
	return result, err
}

// 链接响应消息
type LinkResponseMessage struct {
	XMLName xml.Name `xml:"xml"`
	baseMessage
	MsgId       int       `xml:"MsgId"`
	Title       CDATAText `xml:"Title"`
	Description CDATAText `xml:"Description"`
	Url         CDATAText `xml:"Url"`
}

func NewLinkResponseMessage(toUserName string, fromUserName string, createTime time.Duration, title string, description string, url string, msgId int) LinkResponseMessage {
	return LinkResponseMessage{baseMessage: baseMessage{ToUserName: NewCDATAText(toUserName), FromUserName: NewCDATAText(fromUserName), CreateTime: createTime, MsgType: NewCDATAText("link")}, MsgId: msgId, Title: NewCDATAText(title), Description: NewCDATAText(description), Url: NewCDATAText(url)}
}

func EncryptLinkMsg(config WXAPIConfig, replyMsg LinkResponseMessage, timestamp string, nonce string) (string, error) {
	bReplyMsg, _ := xml.Marshal(&replyMsg)
	sReplyMsg := string(bReplyMsg)
	log.Println("raw_message=" + sReplyMsg)
	result, err := EncryptMsg(config, sReplyMsg, timestamp, nonce)
	log.Println("encrypt_message=" + result)
	return result, err
}

// 事件响应消息
type EventResponseMessage struct {
	XMLName xml.Name `xml:"xml"`
	baseMessage
	Event       CDATAText `xml:"Event"`
	Description CDATAText `xml:"Description"`
	Url         CDATAText `xml:"Url"`
}
