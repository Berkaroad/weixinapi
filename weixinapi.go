// Weixin API
package weixinapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	WXAPI_AccessServer_Error = "-50001"
)

type ErrorMessage struct {
	Error_Code int    `json:"errcode"`
	Error_Msg  string `json:"errmsg"`
}

type AccessToken struct {
	Access_Token string `json:"access_token"`
	Expires_In   int    `json:"expires_in"`
	createTime   time.Time
}

type IPList struct {
	IP_List []string `json:"ip_list"`
}

type JsApiTicket struct {
	ErrorMessage
	Ticket     string `json:"ticket"`
	Expires_In int    `json:"expires_in"`
	createTime time.Time
}

// 模拟缓存
var cachedAccessToken = AccessToken{}
var cachedJsApiTicket = JsApiTicket{}

// 获取AccessToken
func GetAccessToken(appID string, appSecret string) (AccessToken, error) {
	if int(time.Now().Sub(cachedAccessToken.createTime).Seconds()) >= cachedAccessToken.Expires_In {
		url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + appID + "&secret=" + appSecret
		errMsg := ErrorMessage{}
		if resp, err := http.Get(url); err == nil {
			if buffer, err := ioutil.ReadAll(resp.Body); err == nil {
				fmt.Println("weixinapi.GetAccessToken=" + string(buffer))
				json.Unmarshal(buffer, &cachedAccessToken)
				if cachedAccessToken.Access_Token == "" {
					json.Unmarshal(buffer, &errMsg)
					return cachedAccessToken, errors.New(errMsg.Error_Msg)
				} else {
					cachedAccessToken.createTime = time.Now()
					return cachedAccessToken, nil
				}
			}
		}
		return cachedAccessToken, errors.New(WXAPI_AccessServer_Error)
	} else {
		return cachedAccessToken, nil
	}
}

// 获取微信服务器IP列表
func GetCallbackIP(accessToken AccessToken) (IPList, error) {
	url := "https://api.weixin.qq.com/cgi-bin/getcallbackip?access_token=" + accessToken.Access_Token
	ipList := IPList{}
	errMsg := ErrorMessage{}
	if resp, err := http.Get(url); err == nil {
		if buffer, err := ioutil.ReadAll(resp.Body); err == nil {
			fmt.Println("weixinapi.GetCallbackIP=" + string(buffer))
			json.Unmarshal(buffer, &ipList)
			if len(ipList.IP_List) == 0 {
				json.Unmarshal(buffer, &errMsg)
				return ipList, errors.New(errMsg.Error_Msg)
			} else {
				return ipList, nil
			}
		}
	}
	return ipList, errors.New(WXAPI_AccessServer_Error)
}

// 获取JSAPI票据
func GetJsApiTicket(accessToken AccessToken) (JsApiTicket, error) {
	if int(time.Now().Sub(cachedJsApiTicket.createTime).Seconds()) >= cachedJsApiTicket.Expires_In {
		url := "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=" + accessToken.Access_Token + "&type=jsapi"
		if resp, err := http.Get(url); err == nil {
			if buffer, err := ioutil.ReadAll(resp.Body); err == nil {
				fmt.Println("weixinapi.GetJsApiTicket=" + string(buffer))
				json.Unmarshal(buffer, &cachedJsApiTicket)
				if cachedJsApiTicket.Error_Code != 0 {
					return cachedJsApiTicket, errors.New(cachedJsApiTicket.Error_Msg)
				} else {
					cachedJsApiTicket.createTime = time.Now()
					return cachedJsApiTicket, nil
				}
			}
		}
		return cachedJsApiTicket, errors.New(WXAPI_AccessServer_Error)
	} else {
		return cachedJsApiTicket, nil
	}
}
