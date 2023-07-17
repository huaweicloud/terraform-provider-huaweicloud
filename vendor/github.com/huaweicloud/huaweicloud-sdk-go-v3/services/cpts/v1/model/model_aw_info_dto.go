package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AwInfoDto struct {

	// aw名
	Name *string `json:"name,omitempty"`

	// 数据库中dc_case_aw表中的主键ID
	Id *string `json:"id,omitempty"`

	// 数据类型（用例/aw/事务）
	DatumType *int32 `json:"datumType,omitempty"`
}

func (o AwInfoDto) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AwInfoDto struct{}"
	}

	return strings.Join([]string{"AwInfoDto", string(data)}, " ")
}
