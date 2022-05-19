package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateMbTasksReportRequest struct {
	Body *MbTasksReportReq `json:"body,omitempty"`
}

func (o CreateMbTasksReportRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateMbTasksReportRequest struct{}"
	}

	return strings.Join([]string{"CreateMbTasksReportRequest", string(data)}, " ")
}
