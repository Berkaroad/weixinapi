package weixinapi

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"log"
	"sort"
	"strings"
	"time"
)

const (
	WXBizMsgCrypt_OK                      = "0"
	WXBizMsgCrypt_ValidateSignature_Error = "-40001"
	WXBizMsgCrypt_ParseXml_Error          = "-40002"
	WXBizMsgCrypt_ComputeSignature_Error  = "-40003"
	WXBizMsgCrypt_IllegalAesKey           = "-40004"
	WXBizMsgCrypt_ValidateAppid_Error     = "-40005"
	WXBizMsgCrypt_EncryptAES_Error        = "-40006"
	WXBizMsgCrypt_DecryptAES_Error        = "-40007"
	WXBizMsgCrypt_IllegalBuffer           = "-40008"
	WXBizMsgCrypt_EncodeBase64_Error      = "-40009"
	WXBizMsgCrypt_DecodeBase64_Error      = "-40010"
)

const (
	MsgType_Text     = "text"
	MsgType_Location = "location"
	MsgType_Link     = "link"
)

// 加密的请求消息
type EncryptRequestMessage struct {
	XMLName    xml.Name `xml:"xml"`
	ToUserName string   `xml:"ToUserName"`
	Encrypt    string   `xml:"Encrypt"`
}

// 请求消息
type RequestMessage struct {
	XMLName      xml.Name      `xml:"xml"`
	URL          string        `xml:"URL"`
	ToUserName   string        `xml:"ToUserName"`
	FromUserName string        `xml:"FromUserName"`
	CreateTime   time.Duration `xml:"CreateTime"`
	MsgType      string        `xml:"MsgType"`
	Content      string        `xml:"Content"`
	Location_X   float64       `xml:"Location_X"`
	Location_Y   float64       `xml:"Location_Y"`
	Scale        int           `xml:"Scale"`
	Label        string        `xml:"Label"`
	MsgId        int           `xml:"MsgId"`
}

// 加密的响应消息
type EncryptResponseMessage struct {
	XMLName      xml.Name  `xml:"xml"`
	Encrypt      CDATAText `xml:"Encrypt"`
	MsgSignature CDATAText `xml:"MsgSignature"`
	TimeStamp    CDATAText `xml:"TimeStamp"`
	Nonce        CDATAText `xml:"Nonce"`
}

func NewEncryptResponseMessage(encrypt string, msgSignature string, timestamp string, nonce string) EncryptResponseMessage {
	instance := EncryptResponseMessage{Encrypt: NewCDATAText(encrypt), MsgSignature: NewCDATAText(msgSignature), TimeStamp: NewCDATAText(timestamp), Nonce: NewCDATAText(nonce)}
	return instance
}

// Verify Signature
func VerifySignature(signature string, token string, timestamp string, nonce string) bool {
	tmpArr := []string{token, timestamp, nonce}
	sort.Strings(tmpArr)
	tmpStr := strings.Join(tmpArr, "")
	cryptor := sha1.New()
	cryptor.Write([]byte(tmpStr))
	tmpStr = hex.EncodeToString(cryptor.Sum(nil))
	if tmpStr == signature {
		return true
	} else {
		return false
	}
}

// Verify Message Signature
func VerifyMsgSignature(signature string, token string, timestamp string, nonce string, msgEncrypt string) bool {
	tmpArr := []string{token, timestamp, nonce, msgEncrypt}
	sort.Strings(tmpArr)
	tmpStr := strings.Join(tmpArr, "")
	cryptor := sha1.New()
	cryptor.Write([]byte(tmpStr))
	tmpStr = hex.EncodeToString(cryptor.Sum(nil))
	if tmpStr == signature {
		return true
	} else {
		return false
	}
}

// Genarate Message Signature
func GenarateMsgSignature(token string, timestamp string, nonce string, msgEncrypt string) (string, error) {
	tmpArr := []string{token, timestamp, nonce, msgEncrypt}
	sort.Strings(tmpArr)
	tmpStr := strings.Join(tmpArr, "")
	hashcode := ""
	defer func() {
		if r := recover(); r != nil {
			log.Println("GenarateMsgSignature error! code:", WXBizMsgCrypt_ComputeSignature_Error)
		}
	}()
	cryptor := sha1.New()
	cryptor.Write([]byte(tmpStr))
	hashcode = hex.EncodeToString(cryptor.Sum(nil))
	hashcode = strings.ToLower(hashcode)
	return hashcode, nil
}

// Encrypt Message
func EncryptMsg(config WXAPIConfig, replyMsg string, timestamp string, nonce string) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("EncryptMsg error! code:", WXBizMsgCrypt_EncryptAES_Error)
		}
	}()
	key, _ := base64.StdEncoding.DecodeString(config.EncodingAESKey + "=")
	random_str := GenerateRandomString(16)
	bRand := []byte(random_str)
	bAppid := []byte(config.AppId)
	bReplyMsg := []byte(replyMsg)
	bMsgLen := Int32ToByteArray(HostToNetworkOrderInt32(int32(len(bReplyMsg))))
	bMsg := make([]byte, len(bRand)+len(bMsgLen)+len(bAppid)+len(bReplyMsg))
	copy(bMsg, bRand)
	copy(bMsg[len(bRand):], bMsgLen)
	copy(bMsg[len(bRand)+len(bMsgLen):], bReplyMsg)
	copy(bMsg[len(bRand)+len(bMsgLen)+len(bReplyMsg):], bAppid)
	pad_length := 32 - len(bMsg)%32
	text_buffer := make([]byte, len(bMsg)+pad_length)
	copy(text_buffer, bMsg)
	pad := KCS7Encoder(len(bMsg))
	copy(text_buffer[len(bMsg):], pad)

	encrypt_buffer := make([]byte, len(text_buffer))

	block, _ := aes.NewCipher(key)
	mode := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])
	mode.CryptBlocks(encrypt_buffer, text_buffer)
	if timestamp == "" {
		timestamp = GenerateTimestampString()
	}

	encryptMsgBody := base64.StdEncoding.EncodeToString(encrypt_buffer)
	// 生成安全签名
	msgSignature, _ := GenarateMsgSignature(config.Token, timestamp, nonce, encryptMsgBody)
	encryptMsg := NewEncryptResponseMessage(encryptMsgBody, msgSignature, timestamp, nonce)
	bEncryptMsg, _ := xml.Marshal(&encryptMsg)
	return string(bEncryptMsg), nil
}

// Descrypt Message
func DecryptMsg(config WXAPIConfig, msgSignature string, timestamp string, nonce string, encryptMsg EncryptRequestMessage) (RequestMessage, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("DecryptMsg error! code:", WXBizMsgCrypt_DecryptAES_Error)
		}
	}()

	key, _ := base64.StdEncoding.DecodeString(config.EncodingAESKey + "=")
	if isValid := VerifyMsgSignature(msgSignature, config.Token, timestamp, nonce, encryptMsg.Encrypt); isValid {
		encryptMsgBody_buffer, _ := base64.StdEncoding.DecodeString(encryptMsg.Encrypt)
		plain_buffer := make([]byte, len(encryptMsgBody_buffer))
		block, _ := aes.NewCipher(key)
		mode := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])
		mode.CryptBlocks(plain_buffer, encryptMsgBody_buffer)
		pad := int(plain_buffer[len(plain_buffer)-1])
		if pad < 1 || pad > 32 {
			pad = 0
		}
		plain_buffer = plain_buffer[:len(plain_buffer)-pad]
		message_len := int(NetworkToHostOrderInt32(ByteArrayToInt32(plain_buffer, 16)))
		bMsg := make([]byte, message_len)
		bAppid := make([]byte, len(plain_buffer)-20-message_len)
		copy(bMsg, plain_buffer[20:message_len+20])
		copy(bAppid, plain_buffer[20+message_len:])
		if string(bAppid) == config.AppId {
			plainMsg := RequestMessage{}
			xml.Unmarshal(bMsg, &plainMsg)
			return plainMsg, nil
		} else {
			return RequestMessage{}, errors.New(WXBizMsgCrypt_ValidateAppid_Error)
		}
	} else {
		return RequestMessage{}, errors.New(WXBizMsgCrypt_ValidateSignature_Error)
	}
}
