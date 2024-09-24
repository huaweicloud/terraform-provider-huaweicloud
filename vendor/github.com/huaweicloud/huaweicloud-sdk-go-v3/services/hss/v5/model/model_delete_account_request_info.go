package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteAccountRequestInfo 删除账号请求
type DeleteAccountRequestInfo struct {

	// 组织Id
	OrganizationId string `json:"organization_id"`

	// 账号ID
	AccountId string `json:"account_id"`

	// 租户项目ID
	ProjectId string `json:"project_id"`
}

func (o DeleteAccountRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteAccountRequestInfo struct{}"
	}

	return strings.Join([]string{"DeleteAccountRequestInfo", string(data)}, " ")
}
