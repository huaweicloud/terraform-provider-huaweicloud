package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CaseDoc 文本用例信息
type CaseDoc struct {

	// 用例描述信息
	Description *string `json:"description,omitempty"`

	// 标签
	LabelNames *[]string `json:"label_names,omitempty"`

	// 前置条件
	Preparation *string `json:"preparation,omitempty"`

	// 用例等级（0-L0；1-L1；2-L2；3-L3；4-L4；）
	Rank *int32 `json:"rank,omitempty"`

	// 状态（0-新建；5-设计中；6-测试中；7-完成；）
	StatusCode *int32 `json:"status_code,omitempty"`

	// 测试步骤
	Steps *[]CaseDocSteps `json:"steps,omitempty"`
}

func (o CaseDoc) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CaseDoc struct{}"
	}

	return strings.Join([]string{"CaseDoc", string(data)}, " ")
}
