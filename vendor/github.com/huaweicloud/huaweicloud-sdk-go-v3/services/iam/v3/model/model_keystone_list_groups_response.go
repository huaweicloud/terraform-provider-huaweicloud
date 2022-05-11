package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneListGroupsResponse struct {

	// 用户组信息列表。
	Groups *[]KeystoneGroupResult `json:"groups,omitempty"`

	Links          *Links `json:"links,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o KeystoneListGroupsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneListGroupsResponse struct{}"
	}

	return strings.Join([]string{"KeystoneListGroupsResponse", string(data)}, " ")
}
