package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchAddAccountsResponse Response Object
type BatchAddAccountsResponse struct {

	// 批量添加账号结果   - true ：成功   - false ：失败
	IsAllLegalCount *bool `json:"is_all_legal_count,omitempty"`

	XRequestId     *string `json:"X-request-id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o BatchAddAccountsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchAddAccountsResponse struct{}"
	}

	return strings.Join([]string{"BatchAddAccountsResponse", string(data)}, " ")
}
