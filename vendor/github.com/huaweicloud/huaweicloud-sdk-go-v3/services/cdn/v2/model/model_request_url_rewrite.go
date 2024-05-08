package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RequestUrlRewrite 访问URL重写。
type RequestUrlRewrite struct {
	Condition *UrlRewriteCondition `json:"condition"`

	// 重定向状态码。支持301、302、303、307。
	RedirectStatusCode *int32 `json:"redirect_status_code,omitempty"`

	// 重定向URL。重定向后的URL，以正斜线（/）开头，不含http://头及域名，如：/test/index.html。   - 当匹配类型为全路径时，\"\\*\"可以用“$1”捕获，例如：匹配内容为/test/\\*.jpg，重定向URL配置为/newtest/$1.jpg，则用户请求/test/11.jpg时，$1捕获11，重定向后请求的URL为/newtest/11.jpg。
	RedirectUrl string `json:"redirect_url"`

	// 支持将客户端请求重定向到其他域名。   > 不填时默认为当前域名。   > 支持字符长度为1-255，必须以http://或https://开头，例如http://www.example.com。
	RedirectHost *string `json:"redirect_host,omitempty"`

	// 执行规则：   - redirect：如果请求的URL匹配了当前规则，该请求将被重定向到目标Path。执行完当前规则后，当存在其他配置规则时，会继续匹配剩余规则。   - break：如果请求的URL匹配了当前规则，请求将被改写为目标Path。执行完当前规则后，当存在其他配置规则时，将不再匹配剩余规则，此时不支持配置重定向Host和重定向状态码，返回状态码200。
	ExecutionMode string `json:"execution_mode"`
}

func (o RequestUrlRewrite) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RequestUrlRewrite struct{}"
	}

	return strings.Join([]string{"RequestUrlRewrite", string(data)}, " ")
}
