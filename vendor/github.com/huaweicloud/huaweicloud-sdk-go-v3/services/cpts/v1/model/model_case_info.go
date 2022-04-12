package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CaseInfo struct {
	// case_id

	CaseId *int32 `json:"case_id,omitempty"`
	// case_name

	CaseName *string `json:"case_name,omitempty"`
	// case_type

	CaseType *int32 `json:"case_type,omitempty"`
	// contents

	Contents *[]Contents `json:"contents,omitempty"`
	// for_loop_params

	ForLoopParams *[]interface{} `json:"for_loop_params,omitempty"`
	// increase_setting

	IncreaseSetting *[]interface{} `json:"increase_setting,omitempty"`
	// stages

	Stages *[]interface{} `json:"stages,omitempty"`
	// status

	Status *int32 `json:"status,omitempty"`
	// temp_id

	TempId *int32 `json:"temp_id,omitempty"`
	// sort

	Sort *int32 `json:"sort,omitempty"`
}

func (o CaseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CaseInfo struct{}"
	}

	return strings.Join([]string{"CaseInfo", string(data)}, " ")
}
