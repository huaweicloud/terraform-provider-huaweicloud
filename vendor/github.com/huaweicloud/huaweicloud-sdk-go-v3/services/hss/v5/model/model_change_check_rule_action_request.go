package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ChangeCheckRuleActionRequest Request Object
type ChangeCheckRuleActionRequest struct {

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 主机ID，不赋值时，查租户所有主机
	HostId *string `json:"host_id,omitempty"`

	// 动作 - \"ignore\" - \"unignore\" - \"fix\" - \"verify\"
	Action string `json:"action"`

	Body *CheckRuleIdListRequestInfo `json:"body,omitempty"`
}

func (o ChangeCheckRuleActionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeCheckRuleActionRequest struct{}"
	}

	return strings.Join([]string{"ChangeCheckRuleActionRequest", string(data)}, " ")
}
