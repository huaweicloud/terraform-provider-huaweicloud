package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListDevicePoliciesResponse Response Object
type ListDevicePoliciesResponse struct {

	// 策略信息列表。
	Policies *[]ListDevicePolicyBase `json:"policies,omitempty"`

	Page           *Page `json:"page,omitempty"`
	HttpStatusCode int   `json:"-"`
}

func (o ListDevicePoliciesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListDevicePoliciesResponse struct{}"
	}

	return strings.Join([]string{"ListDevicePoliciesResponse", string(data)}, " ")
}
