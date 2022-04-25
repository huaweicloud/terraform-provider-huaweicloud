package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneListEndpointsResponse struct {
	Links *Links `json:"links,omitempty"`

	// 终端节点信息列表。
	Endpoints      *[]Endpoint `json:"endpoints,omitempty"`
	HttpStatusCode int         `json:"-"`
}

func (o KeystoneListEndpointsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneListEndpointsResponse struct{}"
	}

	return strings.Join([]string{"KeystoneListEndpointsResponse", string(data)}, " ")
}
