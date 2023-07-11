package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateNewCaseRequest Request Object
type CreateNewCaseRequest struct {
	Body *CaseInfoDetail `json:"body,omitempty"`
}

func (o CreateNewCaseRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateNewCaseRequest struct{}"
	}

	return strings.Join([]string{"CreateNewCaseRequest", string(data)}, " ")
}
