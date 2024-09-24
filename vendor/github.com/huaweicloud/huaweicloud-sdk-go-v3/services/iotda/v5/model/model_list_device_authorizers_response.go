package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListDeviceAuthorizersResponse Response Object
type ListDeviceAuthorizersResponse struct {

	// **参数说明**：自定义鉴权列表。
	Authorizers *[]DeviceAuthorizerSimple `json:"authorizers,omitempty"`

	Page           *Page `json:"page,omitempty"`
	HttpStatusCode int   `json:"-"`
}

func (o ListDeviceAuthorizersResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListDeviceAuthorizersResponse struct{}"
	}

	return strings.Join([]string{"ListDeviceAuthorizersResponse", string(data)}, " ")
}
