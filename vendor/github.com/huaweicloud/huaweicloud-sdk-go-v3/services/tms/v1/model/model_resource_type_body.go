package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ResourceTypeBody struct {

	// 资源类型
	ResourceType string `json:"resource_type"`

	// 资源类型显示名称，可以通过参数中“locale”设置语言。
	ResourceTypeI18nDisplayName string `json:"resource_type_i18n_display_name"`

	// 支持的region列表
	Regions []string `json:"regions"`

	// 是否是全局类型的资源
	Global bool `json:"global"`
}

func (o ResourceTypeBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResourceTypeBody struct{}"
	}

	return strings.Join([]string{"ResourceTypeBody", string(data)}, " ")
}
