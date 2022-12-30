package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CloseProtectionInfoRequestInfo struct {

	// 需要关闭勒索防护的主机ID列表
	HostIdList *[]string `json:"host_id_list,omitempty"`

	// 需要关闭勒索防护的agentID列表
	AgentIdList *[]string `json:"agent_id_list,omitempty"`

	// 关闭防护类型，包含如下：   - close_all : 关闭所有防护   - close_anti : 关闭勒索防护   - close_backup : 关闭备份功能
	CloseProtectionType *string `json:"close_protection_type,omitempty"`
}

func (o CloseProtectionInfoRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CloseProtectionInfoRequestInfo struct{}"
	}

	return strings.Join([]string{"CloseProtectionInfoRequestInfo", string(data)}, " ")
}
