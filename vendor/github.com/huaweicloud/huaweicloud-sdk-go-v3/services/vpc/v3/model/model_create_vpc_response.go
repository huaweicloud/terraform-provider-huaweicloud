package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateVpcResponse struct {

	// 请求ID
	RequestId *string `json:"request_id,omitempty"`

	// 错误消息
	ErrorMsg *string `json:"error_msg,omitempty"`

	// 错误码
	ErrorCode      *string `json:"error_code,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateVpcResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateVpcResponse struct{}"
	}

	return strings.Join([]string{"CreateVpcResponse", string(data)}, " ")
}
