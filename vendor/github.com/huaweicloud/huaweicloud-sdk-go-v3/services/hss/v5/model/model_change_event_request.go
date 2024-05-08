package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ChangeEventRequest Request Object
type ChangeEventRequest struct {

	// Region ID
	Region string `json:"region"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 容器实例名称
	ContainerName *string `json:"container_name,omitempty"`

	// 容器Id
	ContainerId *string `json:"container_id,omitempty"`

	Body *ChangeEventRequestInfo `json:"body,omitempty"`
}

func (o ChangeEventRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeEventRequest struct{}"
	}

	return strings.Join([]string{"ChangeEventRequest", string(data)}, " ")
}
