package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdktime"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ProjectsSet struct {
	// 创建时间

	CreateTime *sdktime.SdkTime `json:"CreateTime,omitempty"`
	// 更新时间

	UpdateTime *sdktime.SdkTime `json:"UpdateTime,omitempty"`
	// 描述

	Description *string `json:"description,omitempty"`
	// 工程id

	Id *int32 `json:"id,omitempty"`
	// 工程名字

	Name *string `json:"name,omitempty"`
	// 工程状态

	Status *int32 `json:"status,omitempty"`
	// 外部参数

	ExternalParams *interface{} `json:"external_params,omitempty"`
	// 文件变量

	VariablesNoFile *[]string `json:"variables_no_file,omitempty"`
}

func (o ProjectsSet) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProjectsSet struct{}"
	}

	return strings.Join([]string{"ProjectsSet", string(data)}, " ")
}
