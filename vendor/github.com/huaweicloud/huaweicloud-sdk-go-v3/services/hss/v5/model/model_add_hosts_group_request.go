package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddHostsGroupRequest Request Object
type AddHostsGroupRequest struct {

	// Region ID
	Region string `json:"region"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 缺省值:application/json; charset=utf-8
	ContentType *string `json:"Content-Type,omitempty"`

	Body *AddHostsGroupRequestInfo `json:"body,omitempty"`
}

func (o AddHostsGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddHostsGroupRequest struct{}"
	}

	return strings.Join([]string{"AddHostsGroupRequest", string(data)}, " ")
}
