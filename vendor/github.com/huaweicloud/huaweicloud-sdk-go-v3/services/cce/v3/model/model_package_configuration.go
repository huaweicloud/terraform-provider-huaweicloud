package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PackageConfiguration struct {

	// 组件名称
	Name *string `json:"name,omitempty"`

	// 组件配置项
	Configurations *[]ConfigurationItem `json:"configurations,omitempty"`
}

func (o PackageConfiguration) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PackageConfiguration struct{}"
	}

	return strings.Join([]string{"PackageConfiguration", string(data)}, " ")
}
