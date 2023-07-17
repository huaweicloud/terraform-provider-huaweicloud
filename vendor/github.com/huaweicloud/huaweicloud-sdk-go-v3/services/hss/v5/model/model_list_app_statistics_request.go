package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAppStatisticsRequest Request Object
type ListAppStatisticsRequest struct {

	// 软件名称
	AppName *string `json:"app_name,omitempty"`

	// 企业项目
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 默认10
	Limit *int32 `json:"limit,omitempty"`

	// 偏移量，为页数*每页显示条数
	Offset *int32 `json:"offset,omitempty"`
}

func (o ListAppStatisticsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAppStatisticsRequest struct{}"
	}

	return strings.Join([]string{"ListAppStatisticsRequest", string(data)}, " ")
}
