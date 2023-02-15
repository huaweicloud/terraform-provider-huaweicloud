package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ProviderResponseBody struct {

	// 云服务名称
	Provider string `json:"provider"`

	// 服务显示名称，可以通过参数中“locale”设置语言。
	ProviderI18nDisplayName string `json:"provider_i18n_display_name"`

	// 资源类型列表
	ResourceTypes []ResourceTypeBody `json:"resource_types"`
}

func (o ProviderResponseBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProviderResponseBody struct{}"
	}

	return strings.Join([]string{"ProviderResponseBody", string(data)}, " ")
}
