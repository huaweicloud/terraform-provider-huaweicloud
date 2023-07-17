package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DetailItem struct {

	// 数据库中dc_case_aw表中的主键ID
	CaseAwId *string `json:"caseAwId,omitempty"`

	// 数据类型（用例/aw/事务）
	DatumType *int32 `json:"datumType,omitempty"`

	// 用例/aw/事务名
	Name *string `json:"name,omitempty"`

	// 事务ID
	TransactionId *string `json:"transactionId,omitempty"`

	// aw列表
	AwList *[]DetailItem `json:"awList,omitempty"`
}

func (o DetailItem) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DetailItem struct{}"
	}

	return strings.Join([]string{"DetailItem", string(data)}, " ")
}
