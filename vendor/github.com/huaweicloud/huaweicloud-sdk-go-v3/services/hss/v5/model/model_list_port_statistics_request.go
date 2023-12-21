package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListPortStatisticsRequest Request Object
type ListPortStatisticsRequest struct {

	// 端口号，精确匹配
	Port *int32 `json:"port,omitempty"`

	// 端口字符串，用来进行模糊匹配
	PortString *string `json:"port_string,omitempty"`

	// 端口类型
	Type *string `json:"type,omitempty"`

	// 企业项目
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 排序的key值，目前支持按照端口号port排序
	SortKey *string `json:"sort_key,omitempty"`

	// 升序还是降序，默认升序，asc
	SortDir *string `json:"sort_dir,omitempty"`

	// 默认10
	Limit *int32 `json:"limit,omitempty"`

	// 默认是0
	Offset *int32 `json:"offset,omitempty"`
}

func (o ListPortStatisticsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPortStatisticsRequest struct{}"
	}

	return strings.Join([]string{"ListPortStatisticsRequest", string(data)}, " ")
}
