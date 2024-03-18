package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreatePromInstanceRequest Request Object
type CreatePromInstanceRequest struct {

	// Prometheus实例所属Region，一般为承载REST服务端点的服务器域名或IP，不同服务不同区域的名称不同。
	Region string `json:"region"`

	Body *PromInstanceRequestModel `json:"body,omitempty"`
}

func (o CreatePromInstanceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreatePromInstanceRequest struct{}"
	}

	return strings.Join([]string{"CreatePromInstanceRequest", string(data)}, " ")
}
