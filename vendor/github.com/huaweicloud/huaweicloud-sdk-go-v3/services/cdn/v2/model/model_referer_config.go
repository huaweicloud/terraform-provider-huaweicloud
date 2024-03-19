package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RefererConfig 防盗链。
type RefererConfig struct {

	// Referer黑白名单类型，off：关闭Referer黑白名单，black：Referer黑名单，white：Referer白名单。
	Type string `json:"type"`

	// 域名或IP地址，以“,”进行分割，域名、IP地址可以混合输入，支持泛域名和带端口的域名。域名、IP地址总数不超过400个，端口取值范围1-65535。
	Value *string `json:"value,omitempty"`

	// 是否包含空Referer，如果是黑名单并开启该选项，则表示无referer不允许访问，如果是白名单并开启该选项，则表示无referer允许访问，true:包含空，false：不包含空，默认值false。
	IncludeEmpty *bool `json:"include_empty,omitempty"`
}

func (o RefererConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RefererConfig struct{}"
	}

	return strings.Join([]string{"RefererConfig", string(data)}, " ")
}
