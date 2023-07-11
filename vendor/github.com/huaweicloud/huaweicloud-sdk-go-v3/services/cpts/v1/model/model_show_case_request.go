package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowCaseRequest Request Object
type ShowCaseRequest struct {

	// 用例id
	CaseId int32 `json:"case_id"`
}

func (o ShowCaseRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowCaseRequest struct{}"
	}

	return strings.Join([]string{"ShowCaseRequest", string(data)}, " ")
}
