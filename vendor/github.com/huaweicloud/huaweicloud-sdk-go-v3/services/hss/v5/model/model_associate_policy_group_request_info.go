package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AssociatePolicyGroupRequestInfo struct {

	// 部署的目标策略组ID
	TargetPolicyGroupId string `json:"target_policy_group_id"`

	// 是否要对全量主机部署策略，如果为true的话，不需填写host_id_list，如果为false的话，需要填写host_id_list
	OperateAll *bool `json:"operate_all,omitempty"`

	// 需要部署策略组的主机ID列表
	HostIdList *[]string `json:"host_id_list,omitempty"`
}

func (o AssociatePolicyGroupRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AssociatePolicyGroupRequestInfo struct{}"
	}

	return strings.Join([]string{"AssociatePolicyGroupRequestInfo", string(data)}, " ")
}
