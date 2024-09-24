package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListUserStatisticsRequest Request Object
type ListUserStatisticsRequest struct {

	// 账号名称，参考windows文件命名规则，支持字母、数字、下划线、中文，特殊字符!@.-等，不包括中文标点符号
	UserName *string `json:"user_name,omitempty"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 每页显示数量
	Limit *int32 `json:"limit,omitempty"`

	// 偏移量：指定返回记录的开始位置
	Offset *int32 `json:"offset,omitempty"`

	// 类别，默认为host，包含如下： - host：主机 - container：容器
	Category *string `json:"category,omitempty"`
}

func (o ListUserStatisticsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListUserStatisticsRequest struct{}"
	}

	return strings.Join([]string{"ListUserStatisticsRequest", string(data)}, " ")
}
