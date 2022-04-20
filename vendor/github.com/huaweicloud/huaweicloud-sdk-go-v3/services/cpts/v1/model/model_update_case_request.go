package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateCaseRequest struct {
	// 用例id

	CaseId int32 `json:"case_id"`
	// 类型

	Target string `json:"target"`

	Body *UpdateCaseRequestBody `json:"body,omitempty"`
}

func (o UpdateCaseRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateCaseRequest struct{}"
	}

	return strings.Join([]string{"UpdateCaseRequest", string(data)}, " ")
}
