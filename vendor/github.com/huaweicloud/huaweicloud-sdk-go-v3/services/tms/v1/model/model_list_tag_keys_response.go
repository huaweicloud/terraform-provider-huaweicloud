package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListTagKeysResponse Response Object
type ListTagKeysResponse struct {

	// 标签键列表
	Keys *[]string `json:"keys,omitempty"`

	PageInfo       *PageInfoTagKeys `json:"page_info,omitempty"`
	HttpStatusCode int              `json:"-"`
}

func (o ListTagKeysResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTagKeysResponse struct{}"
	}

	return strings.Join([]string{"ListTagKeysResponse", string(data)}, " ")
}
