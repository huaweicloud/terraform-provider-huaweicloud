package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchAssociateKeypairResponse Response Object
type BatchAssociateKeypairResponse struct {

	// 批量绑定密钥对任务。
	Tasks          *[]TaskResponseBody `json:"tasks,omitempty"`
	HttpStatusCode int                 `json:"-"`
}

func (o BatchAssociateKeypairResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchAssociateKeypairResponse struct{}"
	}

	return strings.Join([]string{"BatchAssociateKeypairResponse", string(data)}, " ")
}
