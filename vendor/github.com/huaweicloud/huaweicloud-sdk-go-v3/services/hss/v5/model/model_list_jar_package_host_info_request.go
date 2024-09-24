package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListJarPackageHostInfoRequest Request Object
type ListJarPackageHostInfoRequest struct {

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 文件名称
	FileName string `json:"file_name"`

	// 类别，包含如下:   - host : 主机   - container : 容器
	Category *string `json:"category,omitempty"`

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 服务器IP
	HostIp *string `json:"host_ip,omitempty"`

	// 每页显示数量
	Limit *int32 `json:"limit,omitempty"`

	// 偏移量：指定返回记录的开始位置
	Offset *int32 `json:"offset,omitempty"`

	// 是否模糊匹配，默认false表示精确匹配
	PartMatch *bool `json:"part_match,omitempty"`
}

func (o ListJarPackageHostInfoRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListJarPackageHostInfoRequest struct{}"
	}

	return strings.Join([]string{"ListJarPackageHostInfoRequest", string(data)}, " ")
}
