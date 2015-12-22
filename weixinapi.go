// Weixin API
package weixinapi

const (
	WXAPI_AccessServer_Error = "-50001"
)

type ErrorMessage struct {
	Error_Code int    `json:"errcode"`
	Error_Msg  string `json:"errmsg"`
}

// 微信公众号配置
type WXAPIConfig struct {
	AppId          string
	AppSecret      string
	EncodingAESKey string
	Token          string
}

func NewWXAPIConfig(appId, appSecret, encodingAESKey, token string) WXAPIConfig {
	return WXAPIConfig{AppId: appId, AppSecret: appSecret, EncodingAESKey: encodingAESKey, Token: token}
}
