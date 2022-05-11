package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListPermanentAccessKeysResponse struct {

	// 认证结果信息列表。
	Credentials    *[]Credentials `json:"credentials,omitempty"`
	HttpStatusCode int            `json:"-"`
}

func (o ListPermanentAccessKeysResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPermanentAccessKeysResponse struct{}"
	}

	return strings.Join([]string{"ListPermanentAccessKeysResponse", string(data)}, " ")
}
