package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Sni 回源SNI。
type Sni struct {

	// 是否开启回源SNI，on：打开，off：关闭。
	Status string `json:"status"`

	// CDN节点回源需要访问的源站域名。如test.example.com。   > 1. 开启回源SNI时必填。   > 2. 不支持泛域名，仅支持输入数字、“-”、“.”、英文大小写字符。
	ServerName *string `json:"server_name,omitempty"`
}

func (o Sni) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Sni struct{}"
	}

	return strings.Join([]string{"Sni", string(data)}, " ")
}
