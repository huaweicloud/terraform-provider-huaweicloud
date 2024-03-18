package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type Project struct {

	// 创建时间
	CreateTime *string `json:"create_time,omitempty"`

	// 描述
	Description *string `json:"description,omitempty"`

	// 租户id（domain_id）
	Group *string `json:"group,omitempty"`

	// 工程id
	Id *int32 `json:"id,omitempty"`

	// 工程名称
	Name *string `json:"name,omitempty"`

	// 来源（0-PerfTest；2-CloudTest）
	Source *int32 `json:"source,omitempty"`

	// 更新时间
	UpdateTime *string `json:"update_time,omitempty"`
}

func (o Project) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Project struct{}"
	}

	return strings.Join([]string{"Project", string(data)}, " ")
}
