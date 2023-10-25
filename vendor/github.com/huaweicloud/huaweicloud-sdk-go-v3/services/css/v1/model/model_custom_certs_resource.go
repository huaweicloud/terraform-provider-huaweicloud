package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CustomCertsResource struct {

	// 证书id。
	Id *string `json:"id,omitempty"`

	// 证书名称。
	FileName *string `json:"fileName,omitempty"`

	// 证书路径。
	FileLocation *string `json:"fileLocation,omitempty"`

	// 证书状态。
	Status *string `json:"status,omitempty"`

	// 证书上传时间。
	UpdateAt *string `json:"updateAt,omitempty"`
}

func (o CustomCertsResource) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CustomCertsResource struct{}"
	}

	return strings.Join([]string{"CustomCertsResource", string(data)}, " ")
}
