package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateDirectoryRequestBody ge
type UpdateDirectoryRequestBody struct {

	// 目录名称
	Name string `json:"name"`
}

func (o UpdateDirectoryRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDirectoryRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateDirectoryRequestBody", string(data)}, " ")
}
