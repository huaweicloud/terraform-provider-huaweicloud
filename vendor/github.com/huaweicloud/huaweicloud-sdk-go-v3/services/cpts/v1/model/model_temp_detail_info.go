package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TempDetailInfo struct {
	// description

	Description *string `json:"description,omitempty"`
	// id

	Id *int32 `json:"id,omitempty"`
	// 是否被引用

	IsQuoted *bool `json:"is_quoted,omitempty"`
	// name

	Name *string `json:"name,omitempty"`
	// temp_type

	TempType *int32 `json:"temp_type,omitempty"`
	// update_time

	UpdateTime *string `json:"update_time,omitempty"`
}

func (o TempDetailInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TempDetailInfo struct{}"
	}

	return strings.Join([]string{"TempDetailInfo", string(data)}, " ")
}
