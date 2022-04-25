package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListPermanentAccessKeysRequest struct {

	// 待查询的IAM用户ID。
	UserId *string `json:"user_id,omitempty"`
}

func (o ListPermanentAccessKeysRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPermanentAccessKeysRequest struct{}"
	}

	return strings.Join([]string{"ListPermanentAccessKeysRequest", string(data)}, " ")
}
