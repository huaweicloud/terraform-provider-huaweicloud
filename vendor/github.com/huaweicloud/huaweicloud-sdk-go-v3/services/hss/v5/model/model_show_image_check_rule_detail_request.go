package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowImageCheckRuleDetailRequest Request Object
type ShowImageCheckRuleDetailRequest struct {

	// Region ID
	Region *string `json:"region,omitempty"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 镜像类型，包含如下:   - private_image : 私有镜像仓库   - shared_image : 共享镜像仓库   - local_image : 本地镜像   - instance_image : 企业镜像
	ImageType string `json:"image_type"`

	// 组织名称（没有镜像相关信息时，表示查询所有镜像）
	Namespace *string `json:"namespace,omitempty"`

	// 镜像名称
	ImageName *string `json:"image_name,omitempty"`

	// 镜像版本名称
	ImageVersion *string `json:"image_version,omitempty"`

	// 基线名称
	CheckName string `json:"check_name"`

	// 基线类型
	CheckType string `json:"check_type"`

	// 检查项id
	CheckRuleId string `json:"check_rule_id"`

	// 标准类型，包含如下:   - cn_standard : 等保合规标准   - hw_standard : 华为标准   - qt_standard : 青腾标准
	Standard string `json:"standard"`

	// 企业仓库实例ID，swr共享版无需使用该参数
	InstanceId *string `json:"instance_id,omitempty"`
}

func (o ShowImageCheckRuleDetailRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowImageCheckRuleDetailRequest struct{}"
	}

	return strings.Join([]string{"ShowImageCheckRuleDetailRequest", string(data)}, " ")
}
