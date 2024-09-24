package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchAddAccountsRequestInfo 账号列表
type BatchAddAccountsRequestInfo struct {

	// 账号列表表详情
	DataList *[]AddAccountsRequestInfo `json:"data_list,omitempty"`
}

func (o BatchAddAccountsRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchAddAccountsRequestInfo struct{}"
	}

	return strings.Join([]string{"BatchAddAccountsRequestInfo", string(data)}, " ")
}
