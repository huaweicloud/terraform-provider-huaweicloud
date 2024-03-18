package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TempDetailInfo struct {

	// 描述
	Description *string `json:"description,omitempty"`

	// 事务id
	Id *int32 `json:"id,omitempty"`

	// 是否被引用
	IsQuoted *bool `json:"is_quoted,omitempty"`

	// 事务名称
	Name *string `json:"name,omitempty"`

	// 事务类型（已弃用，兼容性保留）
	TempType *int32 `json:"temp_type,omitempty"`

	// 更新时间
	UpdateTime *string `json:"update_time,omitempty"`
}

func (o TempDetailInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TempDetailInfo struct{}"
	}

	return strings.Join([]string{"TempDetailInfo", string(data)}, " ")
}
