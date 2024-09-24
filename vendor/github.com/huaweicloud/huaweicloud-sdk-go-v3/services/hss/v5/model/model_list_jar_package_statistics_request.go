package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListJarPackageStatisticsRequest Request Object
type ListJarPackageStatisticsRequest struct {

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// jar包名称
	FileName *string `json:"file_name,omitempty"`

	// 类别，包含如下:   - host : 主机   - container : 容器
	Category *string `json:"category,omitempty"`

	// 每页显示数量
	Limit *int32 `json:"limit,omitempty"`

	// 偏移量：指定返回记录的开始位置
	Offset *int32 `json:"offset,omitempty"`
}

func (o ListJarPackageStatisticsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListJarPackageStatisticsRequest struct{}"
	}

	return strings.Join([]string{"ListJarPackageStatisticsRequest", string(data)}, " ")
}
