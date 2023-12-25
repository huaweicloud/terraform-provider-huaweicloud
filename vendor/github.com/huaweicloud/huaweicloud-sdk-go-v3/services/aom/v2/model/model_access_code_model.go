package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AccessCodeModel struct {

	// access_code。
	AccessCode *string `json:"access_code,omitempty"`

	// access_code_id。
	AccessCodeId *string `json:"access_code_id,omitempty"`

	// 创建时间。
	CreateAt *int64 `json:"create_at,omitempty"`

	// 状态。
	Status *string `json:"status,omitempty"`
}

func (o AccessCodeModel) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AccessCodeModel struct{}"
	}

	return strings.Join([]string{"AccessCodeModel", string(data)}, " ")
}
