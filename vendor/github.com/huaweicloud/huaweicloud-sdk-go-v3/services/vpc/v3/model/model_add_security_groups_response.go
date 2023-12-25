package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddSecurityGroupsResponse Response Object
type AddSecurityGroupsResponse struct {

	// 请求ID
	RequestId *string `json:"request_id,omitempty"`

	Port           *Port `json:"port,omitempty"`
	HttpStatusCode int   `json:"-"`
}

func (o AddSecurityGroupsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddSecurityGroupsResponse struct{}"
	}

	return strings.Join([]string{"AddSecurityGroupsResponse", string(data)}, " ")
}
