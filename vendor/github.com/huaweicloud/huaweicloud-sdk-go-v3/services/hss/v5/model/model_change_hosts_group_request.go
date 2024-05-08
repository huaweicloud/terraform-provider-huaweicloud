package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ChangeHostsGroupRequest Request Object
type ChangeHostsGroupRequest struct {

	// Region ID
	Region string `json:"region"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 缺省值:application/json; charset=utf-8
	ContentType *string `json:"Content-Type,omitempty"`

	Body *ChangeHostsGroupRequestInfo `json:"body,omitempty"`
}

func (o ChangeHostsGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeHostsGroupRequest struct{}"
	}

	return strings.Join([]string{"ChangeHostsGroupRequest", string(data)}, " ")
}
