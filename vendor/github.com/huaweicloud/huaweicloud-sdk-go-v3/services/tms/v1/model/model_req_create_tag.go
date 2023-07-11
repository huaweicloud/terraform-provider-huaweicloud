package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ReqCreateTag 创建标签请求
type ReqCreateTag struct {

	// 项目ID，resource_type为region级别服务时为必选项。
	ProjectId *string `json:"project_id,omitempty"`

	// 资源列表
	Resources []ResourceTagBody `json:"resources"`

	// 标签列表
	Tags []CreateTagRequest `json:"tags"`
}

func (o ReqCreateTag) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ReqCreateTag struct{}"
	}

	return strings.Join([]string{"ReqCreateTag", string(data)}, " ")
}
