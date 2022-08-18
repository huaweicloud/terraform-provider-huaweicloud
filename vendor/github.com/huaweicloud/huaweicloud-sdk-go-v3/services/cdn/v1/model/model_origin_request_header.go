package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 回源请求头
type OriginRequestHeader struct {

	// 设置回源请求头参数。格式要求：长度1~64，由数字，大小写字母，中划线-组成。
	Name string `json:"name"`

	// 设置回源请求头参数的值。当为删除动作时，可不填。格式要求：长度1~512。不支持中文，不支持变量配置，如：$client_ip,$remote_port等。
	Value *string `json:"value,omitempty"`

	// 回源请求头设置类型。delete：删除，set：设置。同一个请求头字段只允许删除或者设置。设置：若原始回源请求中不存在该字段，先执行新增再执行设置。
	Action string `json:"action"`
}

func (o OriginRequestHeader) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "OriginRequestHeader struct{}"
	}

	return strings.Join([]string{"OriginRequestHeader", string(data)}, " ")
}
