package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// TestCaseBriefInfo 用例简略信息
type TestCaseBriefInfo struct {

	// 用例id
	Id int32 `json:"id"`

	// 用例名称
	Name string `json:"name"`

	// 状态（0-已删除；1-启用；2-停用；）
	Status int32 `json:"status"`

	// 创建时间
	CreateTime string `json:"create_time"`

	// 更新时间
	UpdateTime string `json:"update_time"`

	// 所属目录id
	ParentId int32 `json:"parent_id"`

	// 类型（1-目录；2-用例；）
	Type int32 `json:"type"`
}

func (o TestCaseBriefInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TestCaseBriefInfo struct{}"
	}

	return strings.Join([]string{"TestCaseBriefInfo", string(data)}, " ")
}
