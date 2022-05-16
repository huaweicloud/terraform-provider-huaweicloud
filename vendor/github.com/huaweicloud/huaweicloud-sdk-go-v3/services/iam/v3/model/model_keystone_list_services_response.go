package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneListServicesResponse struct {

	// 服务信息列表。
	Services *[]Service `json:"services,omitempty"`

	Links          *Links `json:"links,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o KeystoneListServicesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneListServicesResponse struct{}"
	}

	return strings.Join([]string{"KeystoneListServicesResponse", string(data)}, " ")
}
