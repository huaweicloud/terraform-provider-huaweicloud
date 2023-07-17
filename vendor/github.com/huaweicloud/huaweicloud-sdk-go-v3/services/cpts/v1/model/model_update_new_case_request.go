package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateNewCaseRequest Request Object
type UpdateNewCaseRequest struct {

	// 用例id
	CaseId int32 `json:"case_id"`

	Body *CaseInfoDetail `json:"body,omitempty"`
}

func (o UpdateNewCaseRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateNewCaseRequest struct{}"
	}

	return strings.Join([]string{"UpdateNewCaseRequest", string(data)}, " ")
}
