package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CreateDirectoryRequestBody struct {

	// 目录名称
	Name string `json:"name"`
}

func (o CreateDirectoryRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateDirectoryRequestBody struct{}"
	}

	return strings.Join([]string{"CreateDirectoryRequestBody", string(data)}, " ")
}
