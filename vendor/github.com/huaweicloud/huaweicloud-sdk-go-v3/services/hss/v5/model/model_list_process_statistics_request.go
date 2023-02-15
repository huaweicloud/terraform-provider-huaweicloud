package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListProcessStatisticsRequest struct {

	// 路径
	Path *string `json:"path,omitempty"`

	// 企业项目
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 默认10
	Limit *int32 `json:"limit,omitempty"`

	// 默认是0
	Offset *int32 `json:"offset,omitempty"`
}

func (o ListProcessStatisticsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListProcessStatisticsRequest struct{}"
	}

	return strings.Join([]string{"ListProcessStatisticsRequest", string(data)}, " ")
}
