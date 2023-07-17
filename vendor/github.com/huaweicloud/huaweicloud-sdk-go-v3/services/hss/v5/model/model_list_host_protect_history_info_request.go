package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListHostProtectHistoryInfoRequest Request Object
type ListHostProtectHistoryInfoRequest struct {

	// Region Id
	Region string `json:"region"`

	// 企业项目
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// Host Id
	HostId string `json:"host_id"`

	// 起始时间
	StartTime int64 `json:"start_time"`

	// 终止时间
	EndTime int64 `json:"end_time"`

	// limit
	Limit int32 `json:"limit"`

	// offset
	Offset int32 `json:"offset"`
}

func (o ListHostProtectHistoryInfoRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListHostProtectHistoryInfoRequest struct{}"
	}

	return strings.Join([]string{"ListHostProtectHistoryInfoRequest", string(data)}, " ")
}
