package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type InstanceLtsBasicInfoResp struct {

	// 实例id
	Id *string `json:"id,omitempty"`

	// 实例名称
	Name *string `json:"name,omitempty"`

	// 引擎名
	EngineName *string `json:"engine_name,omitempty"`

	// 引擎版本
	EngineVersion *string `json:"engine_version,omitempty"`

	// 引擎分类
	EngineCategory *string `json:"engine_category,omitempty"`

	// 实例状态
	Status *string `json:"status,omitempty"`

	// 企业项目id
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 实例进行中的任务
	Actions *[]string `json:"actions,omitempty"`
}

func (o InstanceLtsBasicInfoResp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "InstanceLtsBasicInfoResp struct{}"
	}

	return strings.Join([]string{"InstanceLtsBasicInfoResp", string(data)}, " ")
}
