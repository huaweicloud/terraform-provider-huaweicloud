package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteNewCaseRequest Request Object
type DeleteNewCaseRequest struct {

	// 用例id
	CaseId int32 `json:"case_id"`
}

func (o DeleteNewCaseRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteNewCaseRequest struct{}"
	}

	return strings.Join([]string{"DeleteNewCaseRequest", string(data)}, " ")
}
