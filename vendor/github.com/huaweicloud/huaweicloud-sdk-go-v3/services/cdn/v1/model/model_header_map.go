package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 设置HTTP头参数。取值：\"Content-Disposition\", \"Content-Language\", \"Access-Control-Allow-Origin\",\"Access-Control-Allow-Methods\", \"Access-Control-Max-Age\", \"Access-Control-Expose-Headers\"。
type HeaderMap struct {

	// 指示回复的内容该以何种形式展示
	ContentDisposition *string `json:"Content-Disposition,omitempty"`

	// 说明访问者希望采用的语言或语言组合
	ContentLanguage *string `json:"Content-Language,omitempty"`

	// 指定了该响应的资源是否被允许与给定的origin共享
	AccessControlAllowOrigin *string `json:"Access-Control-Allow-Origin,omitempty"`

	// 明确了客户端所要访问的资源允许使用的方法或方法列表
	AccessControlAllowMethods *string `json:"Access-Control-Allow-Methods,omitempty"`

	// Access-Control-Allow-Methods 和Access-Control-Allow-Headers 提供的信息可以被缓存多久
	AccessControlMaxAge *string `json:"Access-Control-Max-Age,omitempty"`

	// 列出了哪些首部可以作为响应的一部分暴露给外部
	AccessControlExposeHeaders *string `json:"Access-Control-Expose-Headers,omitempty"`
}

func (o HeaderMap) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HeaderMap struct{}"
	}

	return strings.Join([]string{"HeaderMap", string(data)}, " ")
}
