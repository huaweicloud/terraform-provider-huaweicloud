package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RemoveSecurityGroupsResponse Response Object
type RemoveSecurityGroupsResponse struct {

	// 请求ID
	RequestId *string `json:"request_id,omitempty"`

	Port           *Port `json:"port,omitempty"`
	HttpStatusCode int   `json:"-"`
}

func (o RemoveSecurityGroupsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RemoveSecurityGroupsResponse struct{}"
	}

	return strings.Join([]string{"RemoveSecurityGroupsResponse", string(data)}, " ")
}
