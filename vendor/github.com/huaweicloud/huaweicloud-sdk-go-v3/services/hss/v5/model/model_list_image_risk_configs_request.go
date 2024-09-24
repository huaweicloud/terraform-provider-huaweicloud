package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListImageRiskConfigsRequest Request Object
type ListImageRiskConfigsRequest struct {

	// Region ID
	Region *string `json:"region,omitempty"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 镜像类型，包含如下:   - private_image : 私有镜像仓库   - shared_image : 共享镜像仓库   - local_image : 本地镜像   - instance_image : 企业镜像
	ImageType string `json:"image_type"`

	// 偏移量：指定返回记录的开始位置
	Offset *int32 `json:"offset,omitempty"`

	// 每页显示数量
	Limit *int32 `json:"limit,omitempty"`

	// 组织名称
	Namespace *string `json:"namespace,omitempty"`

	// 镜像名称
	ImageName *string `json:"image_name,omitempty"`

	// 镜像版本名称
	ImageVersion *string `json:"image_version,omitempty"`

	// 基线名称
	CheckName *string `json:"check_name,omitempty"`

	// 风险等级，包含如下:   - Security : 安全   - Low : 低危   - Medium : 中危   - High : 高危
	Severity *string `json:"severity,omitempty"`

	// 标准类型，包含如下:   - cn_standard : 等保合规标准   - hw_standard : 华为标准   - qt_standard : 青腾标准
	Standard *string `json:"standard,omitempty"`

	// 企业仓库实例ID，swr共享版无需使用该参数
	InstanceId *string `json:"instance_id,omitempty"`
}

func (o ListImageRiskConfigsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListImageRiskConfigsRequest struct{}"
	}

	return strings.Join([]string{"ListImageRiskConfigsRequest", string(data)}, " ")
}
