package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AddAccountsRequestInfo struct {

	// 组织Id
	OrganizationId string `json:"organization_id"`

	// 账号ID
	AccountId string `json:"account_id"`

	// 账号名称
	AccountName string `json:"account_name"`
}

func (o AddAccountsRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddAccountsRequestInfo struct{}"
	}

	return strings.Join([]string{"AddAccountsRequestInfo", string(data)}, " ")
}
