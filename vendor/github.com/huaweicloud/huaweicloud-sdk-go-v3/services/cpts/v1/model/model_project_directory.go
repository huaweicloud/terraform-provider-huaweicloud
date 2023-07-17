package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ProjectDirectory 用例列表详情
type ProjectDirectory struct {

	// id
	Id int32 `json:"id"`

	// 名称
	Name string `json:"name"`

	// 状态（0：已删除，1：启用，2：停用）
	Status int32 `json:"status"`

	// 创建时间
	CreateTime string `json:"create_time"`

	// 更新时间
	UpdateTime string `json:"update_time"`

	// 父id
	ParentId int32 `json:"parent_id"`

	// 类型（1：目录，2：用例）
	Type int32 `json:"type"`
}

func (o ProjectDirectory) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProjectDirectory struct{}"
	}

	return strings.Join([]string{"ProjectDirectory", string(data)}, " ")
}
