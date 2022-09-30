package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SystemTemplates struct {

	// 配置文件id。
	Id *string `json:"id,omitempty"`

	// 配置文件名称。
	Name *string `json:"name,omitempty"`

	// 配置文件内容。
	ConfContent *string `json:"confContent,omitempty"`

	// 描述。
	Desc *string `json:"desc,omitempty"`
}

func (o SystemTemplates) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SystemTemplates struct{}"
	}

	return strings.Join([]string{"SystemTemplates", string(data)}, " ")
}
