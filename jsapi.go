package weixinapi

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"strconv"
	"strings"
)

type JsApiConfig struct {
	AppId     string `json:"appId"`
	Timestamp int    `json:"timestamp"`
	NonceStr  string `json:"nonceStr"`
	Signature string `json:"signature"`
}

type WXBizJsApi struct {
	JsApiTicket string
	Url         string
}

func NewWXBizJsApi(jsApiTicket string, url string) WXBizJsApi {
	return WXBizJsApi{JsApiTicket: jsApiTicket, Url: url}
}

func (self *WXBizJsApi) GenerateApiConfig() (JsApiConfig, error) {
	noncestr := GenerateRandomString(16)
	timestamp := GenerateTimestamp()
	tmpstr := "jsapi_ticket=" + self.JsApiTicket + "&noncestr=" + noncestr + "&timestamp=" + strconv.Itoa(timestamp) + "&url=" + self.Url
	log.Println("JsApiConfig_Signature_Raw=" + tmpstr)
	hashcode := ""
	cryptor := sha1.New()
	cryptor.Write([]byte(tmpstr))
	hashcode = hex.EncodeToString(cryptor.Sum(nil))
	hashcode = strings.ToLower(hashcode)
	config := JsApiConfig{Timestamp: timestamp, NonceStr: noncestr, Signature: hashcode}

	return config, nil
}
