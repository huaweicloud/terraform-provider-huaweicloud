package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type Referer struct {

	// Referer类型。取值：0代表不设置Referer过滤；1代表黑名单；2代表白名单。默认取值为0。
	RefererType int32 `json:"referer_type"`

	// 请输入域名或IP地址，以“;”进行分割，域名、IP地址可以混合输入，支持泛域名添加。输入的域名、IP地址总数不超过100个。当设置防盗链时，此项必填。
	RefererList *string `json:"referer_list,omitempty"`

	// 是否包含空Referer。如果是黑名单并开启该选项，则表示无referer不允许访问。如果是白名单并开启该选项，则表示无referer允许访问。默认值false。
	IncludeEmpty *bool `json:"include_empty,omitempty"`
}

func (o Referer) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Referer struct{}"
	}

	return strings.Join([]string{"Referer", string(data)}, " ")
}
