package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type StopHotPipelineRequestBody struct {

	// 配置文件名称。
	Name string `json:"name"`
}

func (o StopHotPipelineRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StopHotPipelineRequestBody struct{}"
	}

	return strings.Join([]string{"StopHotPipelineRequestBody", string(data)}, " ")
}
