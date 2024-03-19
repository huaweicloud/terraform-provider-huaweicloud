package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CaseAwInfo struct {

	// aw信息
	Aw *[]AwInfoDto `json:"aw,omitempty"`

	// aw详细信息列表
	AwList *[]DetailItem `json:"awList,omitempty"`

	// 数据库中dc_case_aw表中的主键ID
	CaseAwId *string `json:"caseAwId,omitempty"`

	// 数据库中dc_testcase表中的case_uri
	CaseUri *string `json:"caseUri,omitempty"`

	// 数据库中dc_testcase表中的case_uri_iteration
	CaseUriI *string `json:"case_uri_i,omitempty"`

	// 数据类型（用例/aw/事务）
	DatumType *int32 `json:"datumType,omitempty"`

	// 数据库中dc_case_aw表中的主键ID
	Id *string `json:"id,omitempty"`

	// 用例名
	Name *string `json:"name,omitempty"`

	// 数据库中dc_testcase表中的testcase_id
	TaskExecId *string `json:"taskExecId,omitempty"`

	// 任务ID
	TaskId *string `json:"taskId,omitempty"`

	// 数据库中dc_testcase表中的testcase_id
	TestcaseId *string `json:"testcaseId,omitempty"`

	// 事务详细信息列表
	TransactionList *[]DetailItem `json:"transactionList,omitempty"`
}

func (o CaseAwInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CaseAwInfo struct{}"
	}

	return strings.Join([]string{"CaseAwInfo", string(data)}, " ")
}
