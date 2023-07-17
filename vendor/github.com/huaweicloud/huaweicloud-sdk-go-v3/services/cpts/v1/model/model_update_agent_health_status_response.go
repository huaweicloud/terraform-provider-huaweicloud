package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateAgentHealthStatusResponse Response Object
type UpdateAgentHealthStatusResponse struct {

	// 响应码
	Code *string `json:"code,omitempty"`

	// 响应消息
	Message *string `json:"message,omitempty"`

	// 扩展信息
	Extend *interface{} `json:"extend,omitempty"`

	Result         *UpdateAgentStatusResponseDetail `json:"result,omitempty"`
	HttpStatusCode int                              `json:"-"`
}

func (o UpdateAgentHealthStatusResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAgentHealthStatusResponse struct{}"
	}

	return strings.Join([]string{"UpdateAgentHealthStatusResponse", string(data)}, " ")
}
