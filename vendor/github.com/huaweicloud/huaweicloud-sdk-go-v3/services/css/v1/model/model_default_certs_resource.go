package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DefaultCertsResource struct {

	// 证书名称。
	FileName *string `json:"fileName,omitempty"`

	// 证书路径。
	FileLocation *string `json:"fileLocation,omitempty"`

	// 证书状态。
	Status *string `json:"status,omitempty"`

	// 描述列。
	Column *string `json:"column,omitempty"`

	// 证书描述。
	Desc *string `json:"desc,omitempty"`
}

func (o DefaultCertsResource) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DefaultCertsResource struct{}"
	}

	return strings.Join([]string{"DefaultCertsResource", string(data)}, " ")
}
